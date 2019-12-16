package consume

import (
	"strconv"
	"time"

	"git.cchntek.com/Cypress/sts/env"
	"git.cchntek.com/Cypress/sts/module/db"
	"git.cchntek.com/Cypress/sts/structs"
	"git.cchntek.com/CypressModule/dbpool"
	"github.com/jinzhu/gorm"
)

type (
	UserData struct {
		ID          int    `gorm:"AUTO_INCREMENT"`
		PID         int    `gorm:"column:parentid"`
		OID         int    `gorm:"column:ownerid"`
		Owner       string `gorm:"-" json:"owner"`
		Parent      string `gorm:"-" json:"parent"`
		UserID      string `gorm:"column:id" json:"userid"`
		Currency    string `gorm:"column:currency" json:"currency"`
		TotalSignUp int    `gorm:"-" `
		Date        string `gorm:"-" json:"-"`
		TargetName  string `gorm:"_"`

		CreateTime time.Time `gorm:"-" json:"date"`
	}
	Users struct {
		List []UserData
	}
)

func NewUserDate() *UserData {
	return &UserData{}
}

func NewUsers() Users {
	return Users{}
}

func (u *UserData) Do() (err error) {

	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return
	}

	if u.TargetName == structs.TableStatisticParent {
		return u.StatisticParent(mdb)
	}

	if u.TargetName == structs.TableStatisticHall {
		return u.StatisticOwner(mdb)
	}

	return
}

func (u *UserData) Preprocess() error {
	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return err
	}

	user, exist := PreFields[u.UserID]
	if exist {
		u.ID = user.UID
		u.OID = user.OID
		u.PID = user.PID

	} else {
		userRecord, err := db.UserRecord(mdb, u.UserID)
		if err != nil {
			return err
		}

		u.ID = userRecord.ID
		u.OID = userRecord.OID
		u.PID = userRecord.PID
	}
	u.Date = u.CreateTime.In(env.TimeZone).Format(StsDateFormat)
	return nil
}

func (us Users) Arrange(length int) map[string]UserData {

	userMap := make(map[string]UserData, length)

	for k := range us.List {
		if us.List[k].PID == 0 {
			continue
		}

		//key PID +date
		//statistic_parent
		pKey := strconv.Itoa(us.List[k].PID) + us.List[k].Date
		if u, exist := userMap[pKey]; exist {
			u.TotalSignUp++
			u.TargetName = structs.TableStatisticParent

			userMap[pKey] = u
		} else {
			userMap[pKey] = UserData{
				PID:         us.List[k].PID,
				Date:        us.List[k].Date,
				TotalSignUp: 1,
				TargetName:  structs.TableStatisticParent,
			}
		}

		if us.List[k].OID == 0 {
			continue
		}

		//key OID + PID +date
		//statistic_owner
		oKey := strconv.Itoa(us.List[k].OID) + strconv.Itoa(us.List[k].PID) + us.List[k].Date
		if u, exist := userMap[oKey]; exist {
			u.TotalSignUp++
			u.TargetName = structs.TableStatisticHall

			userMap[oKey] = u
		} else {
			userMap[oKey] = UserData{
				OID:         us.List[k].OID,
				PID:         us.List[k].PID,
				Date:        us.List[k].Date,
				TotalSignUp: 1,
				TargetName:  structs.TableStatisticHall,
			}
		}
	}

	return userMap
}

//StatisticParent 代理的統計
func (u *UserData) StatisticParent(tx *gorm.DB) (err error) {

	sql := "INSERT INTO " + structs.TableStatisticParent + "( `pid` ,`date`, `total_sign_up`, `total_user` ,`total_bet`,`total_win`,`total_game`, `total_login`, `total_game_login`,`pct_of_comm`) VALUES(?,?,?,0,0,0,0,0,0,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ON DUPLICATE KEY UPDATE `total_sign_up`=`total_sign_up`+?   "
	sqlArg := []interface{}{u.PID, u.Date, u.TotalSignUp, u.PID, u.TotalSignUp}
	return db.ExecRaw(tx, sql, sqlArg)
}

//StatisticOwner 站長的統計 ☆
func (u *UserData) StatisticOwner(tx *gorm.DB) (err error) {

	sql := "INSERT INTO " + structs.TableStatisticHall + "( `oid`,`date`, `total_sign_up`,`total_user`,`total_bet`,`total_win`,`total_game`, `total_login`, `total_game_login`,`pct_of_comm` ) VALUES(?,?,?,0,0,0,0,0,0,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ON DUPLICATE KEY UPDATE `total_sign_up`=`total_sign_up`+?   "
	sqlArg := []interface{}{u.OID, u.Date, u.TotalSignUp, u.OID, u.TotalSignUp}
	return db.ExecRaw(tx, sql, sqlArg)
}
