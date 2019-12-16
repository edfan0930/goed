package consume

import (
	"strconv"
	"strings"
	"time"

	"git.cchntek.com/Cypress/sts/env"
	"git.cchntek.com/Cypress/sts/module/db"
	"git.cchntek.com/Cypress/sts/structs"
	"git.cchntek.com/CypressModule/dbpool"
	"github.com/allegro/bigcache"
	"github.com/jinzhu/gorm"
)

const (
	DB            = "maria"
	StsDateFormat = "2006-01-02 15:00:00"
)

type (
	//Order used for Bind data , row write
	Order struct {
		OrderID    string    `gorm:"column:id" json:"OrderID"`
		GameHall   string    `gorm:"column:gamehall" json:"GameHall"`
		GameType   string    `gorm:"column:gametype" json:"GameType"`
		Platform   string    `gorm:"column:platform" json:"Platform"`
		GameCode   string    `gorm:"column:gamecode" json:"GameCode"`
		Account    string    `gorm:"column:account" json:"Account"`
		Owner      string    `gorm:"column:owner" json:"OwnerID"`
		Parent     string    `gorm:"column:parent" json:"ParentID"`
		PlayerID   string    `gorm:"column:playerid" json:"PlayerID"`
		GameToken  string    `gorm:"column:gametoken" json:"GameToken"`
		Wins       float64   `gorm:"column:total_win" json:"Wins"`
		Bets       float64   `gorm:"column:total_bet" json:"Bets"`
		Jackpots   float64   `gorm:"column:total_jackpot" json:"Jackpots"`
		RoundID    string    `gorm:"-" json:"RoundID"`
		IndexID    string    `gorm:"-" json:"IndexID" bson:"indexid"`
		Round      int       `gorm:"-" `
		CreateTime time.Time `gorm:"column:datetime" json:"CreateTime"`
		Date       string    `gorm:"column:date" `
		UID        int       `gorm:"-" `
		PID        int       `gorm:"-" `
		GID        int       `gorm:"-" `
		OID        int       `gorm:"-" `
		Currency   string    `gorm:"-" `
		TableName  string    `gorm:"-"`
	}
	Orders struct {
		List []Order
	}
	PreField struct {
		UID      int
		OID      int
		PID      int
		GID      int
		Currency string
	}
)

var PreFields = make(map[string]PreField)
var GID = make(map[string]int)

//Order
func NewOrder() *Order {
	return &Order{}
}

func NewOrders() Orders {
	return Orders{}
}

//Do
func (o *Order) Do() (err error) {

	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return
	}

	if o.TableName == structs.TableStatisticUserByGame {
		return o.StsUserByGame(mdb)
	}

	if o.TableName == structs.TableStatisticParentByGame {
		return o.StsParentByGame(mdb)
	}

	if o.TableName == structs.TableStatisticHallByGame {
		return o.StsHallByGame(mdb)
	}

	return
}

//Preprocess
//todo 優化cache
func (o *Order) Preprocess() error {

	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if err != nil {
		return err
	}

	user, exist := PreFields[o.PlayerID]
	if exist {
		o.UID = user.UID
		o.OID = user.OID
		o.PID = user.PID
		o.Currency = user.Currency

	} else {
		u, err := db.UserRecord(mdb, o.PlayerID)
		if err != nil {
			return err
		}
		PreFields[o.PlayerID] = PreField{
			UID:      u.ID,
			OID:      u.OID,
			PID:      u.PID,
			Currency: u.Currency,
		}
		o.UID = u.ID
		o.OID = u.OID
		o.PID = u.PID
		o.Currency = u.Currency
	}

	gid, exist := GID[o.GameCode]
	if exist {
		o.GID = gid
	} else {
		game, err := db.GameRecord(mdb, o.GameCode)
		if err != nil {
			return err
		}
		GID[o.GameCode] = game.ID
		o.GID = game.ID
	}

	o.Date = o.CreateTime.In(env.TimeZone).Format(StsDateFormat)
	return nil
}

//Preprocess2
func (o *Order) Preprocess2() error {
	cacheConfig := bigcache.DefaultConfig(30 * time.Minute)
	cacheConfig.MaxEntriesInWindow = 2000
	cache, err := bigcache.NewBigCache(cacheConfig)
	if err == nil {
		gidByte, err := cache.Get(o.GameCode)
		if err == nil {
			o.GID, _ = strconv.Atoi(string(gidByte))
		}
		playerByte, err := cache.Get(o.PlayerID)
		if err == nil {
			pSlice := strings.Split(string(playerByte), ";")
			o.UID, _ = strconv.Atoi(pSlice[0])
			o.OID, _ = strconv.Atoi(pSlice[1])
			o.PID, _ = strconv.Atoi(pSlice[2])
			o.Currency = pSlice[3]
		}
	}
	mdb, err := dbpool.GetPool().GetMariaDB(DB)
	if o.GID == 0 {
		game, err := db.GameRecord(mdb, o.GameCode)
		if err != nil {
			return err
		}
		cache.Set(o.GameCode, []byte(strconv.Itoa(game.ID)))
		o.GID = game.ID
	}
	if o.UID == 0 {
		user, err := db.UserRecord(mdb, o.PlayerID)
		if err != nil {
			return err
		}
		uid := strconv.Itoa(o.UID)
		oid := strconv.Itoa(o.OID)
		pid := strconv.Itoa(o.PID)
		v := uid + ";" + oid + ";" + pid + ";" + o.Currency
		cache.Set(v, []byte(o.PlayerID))

		o.UID = user.ID
		o.OID = user.OID
		o.PID = user.PID
		o.Currency = user.Currency
	}

	o.Date = o.CreateTime.In(env.TimeZone).Format(StsDateFormat)
	return nil
}

//Arrange 整理和壓縮資料
func (os Orders) Arrange(length int) map[string]Order {
	//orders := Orders{}
	orderMap := make(map[string]Order, length)

	for k := range os.List {
		//user by game
		//key UID + GID +date
		uByGame := strconv.Itoa(os.List[k].UID) + strconv.Itoa(os.List[k].GID) + os.List[k].Date
		if o, exist := orderMap[uByGame]; exist {

			o.Bets += os.List[k].Bets
			o.Wins += os.List[k].Wins
			o.Jackpots += os.List[k].Jackpots
			o.Round++
			o.TableName = structs.TableStatisticUserByGame

			orderMap[uByGame] = o
		} else {

			orderMap[uByGame] = Order{
				GID:       os.List[k].GID,
				UID:       os.List[k].UID,
				Date:      os.List[k].Date,
				Bets:      os.List[k].Bets,
				Wins:      os.List[k].Wins,
				Jackpots:  os.List[k].Jackpots,
				Round:     1,
				TableName: structs.TableStatisticUserByGame,
			}
		}

		//parent by game
		//key PID + date +GID
		pByGame := strconv.Itoa(os.List[k].PID) + os.List[k].Date + strconv.Itoa(os.List[k].GID)
		if o, exist := orderMap[pByGame]; exist {

			o.Bets += os.List[k].Bets
			o.Wins += os.List[k].Wins
			o.Jackpots += os.List[k].Jackpots
			o.Round++
			o.TableName = structs.TableStatisticParentByGame

			orderMap[pByGame] = o
		} else {
			order := os.List[k]
			order.Round = 1
			order.TableName = structs.TableStatisticParentByGame

			orderMap[pByGame] = order
		}

		//Owner by game
		//key date + OID + GID
		oByGame := os.List[k].Date + strconv.Itoa(os.List[k].OID) + strconv.Itoa(os.List[k].GID)
		newBet, _ := env.RateMap.ChangeToOwnerCurrency(os.List[k].Owner, os.List[k].Bets, os.List[k].Currency)
		newWin, _ := env.RateMap.ChangeToOwnerCurrency(os.List[k].Owner, os.List[k].Wins, os.List[k].Currency)
		newJackpot, _ := env.RateMap.ChangeToOwnerCurrency(os.List[k].Owner, os.List[k].Jackpots, os.List[k].Currency)

		if o, exist := orderMap[oByGame]; exist {

			o.Bets += newBet
			o.Wins += newWin
			o.Jackpots += newJackpot
			o.Round++
			o.TableName = structs.TableStatisticHallByGame

			orderMap[oByGame] = o
		} else {

			order := os.List[k]
			order.Round = 1
			order.TableName = structs.TableStatisticHallByGame

			orderMap[oByGame] = order
		}

	}

	return orderMap
}

func (o *Order) StsUserByGame(d *gorm.DB) error {
	sql := "INSERT INTO `" +
		structs.TableStatisticUserByGame +
		"`(`gid`,`uid`,`date`,`total_bet`,`total_win`,`total_round`,`total_jackpot`) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `total_bet`=`total_bet`+? ,`total_win`=`total_win`+? , `total_round`=`total_round`+? ,`total_jackpot`=`total_jackpot`+? "
	sqlArg := []interface{}{o.GID, o.UID, o.Date, o.Bets, o.Wins, o.Round, o.Jackpots, o.Bets, o.Wins, o.Round, o.Jackpots}
	return db.ExecRaw(d, sql, sqlArg)
}

func (o *Order) StsParentByGame(d *gorm.DB) error {
	sql := "INSERT INTO `" +
		structs.TableStatisticParentByGame +
		"`( `pid`,`gid`,`date`,`total_bet`,`total_win`,`total_round`,`total_user`,`total_jackpot`,`pct_of_comm`  ) VALUES(?,?,?,?,?,?,0,?,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ON DUPLICATE KEY UPDATE `total_bet`=`total_bet`+? ,`total_win`=`total_win`+? ,`total_round`=`total_round`+?  ,`total_jackpot`=`total_jackpot`+?"
	sqlArg := []interface{}{o.PID, o.GID, o.Date, o.Bets, o.Wins, o.Round, o.Jackpots, o.PID, o.Bets, o.Wins, o.Round, o.Jackpots}
	return db.ExecRaw(d, sql, sqlArg)
}

func (o *Order) StsHallByGame(d *gorm.DB) error {
	sql := "INSERT INTO `" + structs.TableStatisticHallByGame + "`( `oid`,`gid`,`date`,`total_bet`,`total_win`,`total_round`,`total_user`,`total_jackpot`,`pct_of_comm` ) VALUES(?,?,?,?,?,?,0,?,(SELECT `pct_of_comm` FROM `parent_list` WHERE id=? ) ) ON DUPLICATE KEY UPDATE `total_bet`=`total_bet`+? ,`total_win`=`total_win`+?  ,`total_round`=`total_round`+?  ,`total_jackpot`=`total_jackpot`+?  "
	sqlArg := []interface{}{o.OID, o.GID, o.Date, o.Bets, o.Wins, o.Round, o.Jackpots, o.OID, o.Bets, o.Wins, o.Round, o.Jackpots}
	return db.ExecRaw(d, sql, sqlArg)
}
