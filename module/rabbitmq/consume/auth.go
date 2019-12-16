package consume

import (
	"strconv"
	"strings"
	"time"

	"git.cchntek.com/Cypress/sts/env"
	"git.cchntek.com/Cypress/sts/module/db"
	"git.cchntek.com/Cypress/sts/structs"
	"git.cchntek.com/CypressModule/dbpool"
	"github.com/jinzhu/gorm"
)

const (
	//WEB 平台 web
	WEB = "web"
	//MOBILE 平台 行動裝置
	MOBILE = "mobile"
	//PC 平台 行動裝置
	PC = "pc"

	//IsUserLogin -使用者登入
	IsUserLogin = "UserLogin"
	//IsGameLogin -遊戲登入
	IsGameLogin = "GameLogin"
)

type (
	LoginInfoData struct {
		GID            int       `gorm:"column:gid"`
		UID            int       `gorm:"column:uid"`
		PID            int       `gorm:"-"`
		OID            int       `gorm:"-"`
		Owner          string    `gorm:"-" json:"owner"`
		Parent         string    `gorm:"-" json:"parent"`
		UserID         string    `gorm:"-" json:"userid"`
		Account        string    `gorm:"-" json:"account"`
		GameCode       string    `gorm:"-" json:"gamecode"`
		Date           string    `gorm:"column:date" json:"-"`
		GameHall       string    `gorm:"-" json:"gamehall"`
		Platform       string    `gorm:"column:platform" json:"platform"`
		LoginTimes     int       `gorm:"column:login_times" json:"login_times"`
		LoginGameTimes int       `gorm:"column:login_game_times"`
		Type           string    `gorm:"-" json:"type"`
		PlatWeb        int       `gorm:"-"`
		PlatMobi       int       `gorm:"-"`
		PlatPC         int       `gorm:"-"`
		TotalRound     int       `gorm:"-"`
		TotalUser      int       `gorm:"-"`
		CreateTime     time.Time `gorm:"-" json:"date"`
		TargetName     string    `gorm:"_"`
	}

	LoginInfos struct {
		List []LoginInfoData
	}
)

func NewLoginInfoData() *LoginInfoData {
	return &LoginInfoData{}
}

func NewLoginInfos() LoginInfos {
	return LoginInfos{}
}

//Preprocess ...
//todo 優化cache
func (l *LoginInfoData) Preprocess() error {

	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return err
	}

	switch strings.ToLower(l.Platform) {
	case MOBILE:
		l.PlatMobi = 1
	case WEB:
		l.PlatWeb = 1
	case PC:
		l.PlatPC = 1
	}

	user, exist := PreFields[l.UserID]
	if exist {
		l.UID = user.UID
		l.OID = user.OID
		l.PID = user.PID

	} else {
		u, err := db.UserRecord(mdb, l.UserID)
		if err != nil {
			return err
		}

		l.UID = u.ID
		l.OID = u.OID
		l.PID = u.PID
	}

	if l.Type == IsGameLogin {
		gid, exist := GID[l.GameCode]
		if exist {
			l.GID = gid
		} else {
			game, err := db.GameRecord(mdb, l.GameCode)
			if err != nil {
				return err
			}

			l.GID = game.ID
		}
	}

	l.Date = l.CreateTime.In(env.TimeZone).Format(StsDateFormat)
	return nil
}

//Do ...
func (l *LoginInfoData) Do() (err error) {
	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return
	}

	if l.TargetName == structs.TableStatisticGame {
		return l.StatisticGame(mdb)
	}

	if l.TargetName == structs.TableStatisticUser {
		return l.StatisticUser(mdb)
	}

	if l.TargetName == "statistic_use_bylogin" {
		return l.StatisticUserForUserLogin(mdb)
	}

	if l.TargetName == structs.TableStatisticParent {
		return l.StatisticParent(mdb)
	}

	if l.TargetName == structs.TableStatisticHall {
		return l.StatisticHall(mdb)
	}

	return
}

//Arrange ...
//todo 優化map 管理
func (ls LoginInfos) Arrange(length int) map[string]LoginInfoData {
	authMap := make(map[string]LoginInfoData, length)
	for k := range ls.List {

		//key GID + date
		//statistic_game
		gameKey := "game" + strconv.Itoa(ls.List[k].GID) + ls.List[k].Date
		if l, exist := authMap[gameKey]; exist {

			l.PlatWeb++
			l.PlatMobi++
			l.PlatPC++
			l.TargetName = structs.TableStatisticGame

			authMap[gameKey] = l
		} else {

			authMap[gameKey] = LoginInfoData{
				GID:        ls.List[k].GID,
				Date:       ls.List[k].Date,
				PlatWeb:    ls.List[k].PlatWeb,
				PlatMobi:   ls.List[k].PlatMobi,
				PlatPC:     ls.List[k].PlatPC,
				TargetName: structs.TableStatisticGame,
			}
		}

		//key UID + date
		//statistic_user
		userKey := "user" + strconv.Itoa(ls.List[k].UID) + ls.List[k].Date
		if ls.List[k].Type == IsGameLogin {
			if l, exist := authMap[userKey]; exist {
				l.LoginGameTimes++
				l.TargetName = structs.TableStatisticUser

				authMap[userKey] = l
			} else {
				authMap[userKey] = LoginInfoData{
					UID:            ls.List[k].UID,
					Date:           ls.List[k].Date,
					LoginTimes:     0,
					LoginGameTimes: 1,
					TargetName:     structs.TableStatisticUser,
				}
			}
		}
		//key UID +date
		//statistic_user
		userLoginKey := "userlogin" + strconv.Itoa(ls.List[k].UID) + ls.List[k].Date
		if ls.List[k].Type == IsUserLogin {
			if l, exist := authMap[userLoginKey]; exist {
				l.LoginTimes++
				l.TargetName = "statistic_use_bylogin"

				authMap[userLoginKey] = l
			} else {
				authMap[userLoginKey] = LoginInfoData{
					UID:            ls.List[k].UID,
					Date:           ls.List[k].Date,
					LoginTimes:     1,
					LoginGameTimes: 0,
					TargetName:     "statistic_use_bylogin",
				}
			}
		}

		//key PID +date
		//statistic_parent
		parentKey := "parent" + strconv.Itoa(ls.List[k].PID) + ls.List[k].Date
		if l, exist := authMap[parentKey]; exist {
			l.TotalUser++
			l.TargetName = structs.TableStatisticParent

			authMap[parentKey] = l
		} else {
			authMap[parentKey] = LoginInfoData{
				PID:        ls.List[k].PID,
				Date:       ls.List[k].Date,
				TotalUser:  1,
				TargetName: structs.TableStatisticParent,
			}
		}

		//key OID +date
		//statistic_owner
		ownerKey := "owner" + strconv.Itoa(ls.List[k].OID) + ls.List[k].Date
		if l, exist := authMap[ownerKey]; exist {
			l.TotalUser++
			l.TargetName = structs.TableStatisticHall

			authMap[ownerKey] = l
		} else {
			authMap[ownerKey] = LoginInfoData{
				OID:        ls.List[k].OID,
				Date:       ls.List[k].Date,
				TotalUser:  1,
				TargetName: structs.TableStatisticHall,
			}
		}
	}

	return authMap
}

//StatisticGame 遊戲統計
func (l *LoginInfoData) StatisticGame(d *gorm.DB) (err error) {

	sql := " INSERT INTO " + structs.TableStatisticGame + "( `gid`, `date`, `plat_web`, `plat_mobile`, `plat_pc`, `total_user`, `total_round`, `total_bets`, `total_wins`) VALUES(?,?,?,?,?,(SELECT COUNT(1) FROM `statistic_user_by_game` WHERE gid=? and date=? ) ,0,0,0 ) ON DUPLICATE KEY UPDATE `plat_web`=`plat_web`+? ,`plat_mobile`=`plat_mobile`+? ,`plat_pc`=`plat_pc`+? ,`total_user`=(SELECT COUNT(1) FROM `statistic_user_by_game` WHERE gid=? and date=? ) ; "
	sqlArg := []interface{}{l.GID, l.Date, l.PlatWeb, l.PlatMobi, l.PlatPC, l.GID, l.Date, l.PlatWeb, l.PlatMobi, l.PlatPC, l.GID, l.Date}
	return db.ExecRaw(d, sql, sqlArg)

}

//StatisticUser 玩家的遊戲統計 for game
func (l *LoginInfoData) StatisticUser(d *gorm.DB) (err error) {

	sql := " INSERT INTO `" + structs.TableStatisticUser + "`( `uid`, `date`, `total_bet`, `total_win`, `login_times`, `login_game_times`) VALUES(?,?,0,0,0,?) ON DUPLICATE KEY UPDATE `login_game_times`=`login_game_times`+? ; "
	sqlArg := []interface{}{l.UID, l.Date, l.LoginGameTimes, l.LoginGameTimes}
	return db.ExecRaw(d, sql, sqlArg)

}

//StatisticUserForUserLogin 玩家的遊戲統計 for ForUserLogin
func (l *LoginInfoData) StatisticUserForUserLogin(d *gorm.DB) (err error) {

	sql := "INSERT INTO `" + structs.TableStatisticUser + "`( `uid`, `date`, `total_bet`, `total_win`, `login_times`, `login_game_times`) VALUES(?,?,0,0,?,0 ) ON DUPLICATE KEY UPDATE `login_times`=`login_times`+? ; "
	sqlArg := []interface{}{l.UID, l.Date, l.LoginTimes, l.LoginTimes}
	return db.ExecRaw(d, sql, sqlArg)

}

//StatisticParent 代理的統計 ☆
//Todo 優化語法
func (l *LoginInfoData) StatisticParent(d *gorm.DB) (err error) {

	if db.CheckDuplicate(d, structs.TableStatisticParent, "pid=? and date=?", []interface{}{l.PID, l.Date}) {
		//有資料 update
		sql := "update statistic_parent par JOIN (select ul.parentid,su.date,count(1) total_user,sum(su.login_game_times) total_game_login,sum(su.login_times) total_login from user_list ul JOIN statistic_user su ON su.uid=ul.id AND ul.parentid=? AND su.date=?) tmp on par.pid=tmp.parentid and par.date=tmp.date SET par.total_user=tmp.total_user,  par.total_game_login=tmp.total_game_login,  par.total_login=tmp.total_login;"
		sqlArg := []interface{}{l.PID, l.Date}
		return db.ExecRaw(d, sql, sqlArg)

	}
	//insert
	sql := " INSERT INTO `" + structs.TableStatisticParent + "`( `pid` ,`date`, `total_user`,`total_bet`,`total_win`,`total_game`, `total_login`, `total_game_login`, `total_sign_up`,`pct_of_comm` ) VALUES(?,?,?,0,0,0,0,0,0,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ; "
	sqlArg := []interface{}{l.PID, l.Date, l.TotalUser, l.PID}
	return db.ExecRaw(d, sql, sqlArg)

}

//StatisticHall 站長的統計 ☆
//Todo 優化語法
func (l *LoginInfoData) StatisticHall(d *gorm.DB) (err error) {

	if db.CheckDuplicate(d, structs.TableStatisticHall, "oid=? and date=?", []interface{}{l.OID, l.Date}) {
		//有資料 update
		sql := "UPDATE statistic_owner own JOIN (SELECT ul.ownerid , su.date, COUNT(1) total_user,SUM(su.login_game_times) total_game_login,SUM(su.login_times) total_login FROM user_list ul JOIN statistic_user su ON su.uid=ul.id AND ul.ownerid=? AND su.date=?) tmp ON own.oid=tmp.ownerid AND own.date=tmp.date SET own.total_user=tmp.total_user, own.total_game_login=tmp.total_game_login, own.total_login=tmp.total_login;"

		sqlArg := []interface{}{l.OID, l.Date}
		return db.ExecRaw(d, sql, sqlArg)
	}

	sql := " INSERT INTO `" + structs.TableStatisticHall + "`( `oid`,`date`,`total_user`,`total_bet`,`total_win`,`total_game`, `total_login`, `total_game_login`, `total_sign_up` ,`pct_of_comm` ) VALUES(?,?,(SELECT COUNT(1) FROM `statistic_user` LEFT JOIN `user_list` ON `statistic_user`.uid=`user_list`.id  WHERE `user_list`.ownerid=? AND `statistic_user`.date=? ),0,0,0,0,0,0,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ;"
	sqlArg := []interface{}{l.OID, l.Date, l.OID, l.Date, l.OID}

	return db.ExecRaw(d, sql, sqlArg)

}
