package mongodb

import (
	"time"

	"github.com/globalsign/mgo"
)

var timeZone = time.FixedZone("UTC-4", -4*60*60)

//MgoInfo Control session
type (
	MgoInfo struct {
		//session 不暴露
		session *mgo.Session

		//Unordered Ordered Bulk語法 , 當有某筆錯誤時 , 是否繼續執行
		Unordered bool

		PoolTimeout time.Duration

		//Mode changes the consistency mode for the session
		Mode mgo.Mode
	}
)

//SetTimeZone ...
func SetTimeZone(tz *time.Location) {
	timeZone = tz
}

//CreateMgoInfo 建立session
func CreateMgoInfo(
	host []string,
	timeout, poolTimeout time.Duration,
	db, user, pw, source string,
	poolLimit, minPollSize, maxIdleTimeMS int) (m *MgoInfo, err error) {

	session, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:         host,
			Timeout:       timeout,
			Database:      db,
			Source:        source,
			Username:      user,
			Password:      pw,
			PoolLimit:     poolLimit,
			MinPoolSize:   minPollSize,
			MaxIdleTimeMS: maxIdleTimeMS,
			FailFast:      true,
		},
	)

	if err != nil {
		return nil, err
	}

	session.SetPoolTimeout(poolTimeout)

	return &MgoInfo{
		PoolTimeout: poolTimeout,
		session:     session}, nil
}

//NewMgoInfo copy session
//set pool timeout
/* func NewMgoInfo(s *mgo.Session) *MgoInfo {

	return &MgoInfo{
		session: s.Copy(),
	}
} */

//Close close session
func (m *MgoInfo) Close() {
	m.session.Close()
}

//Copy copy session
func (m *MgoInfo) Copy() *MgoInfo {

	return &MgoInfo{
		session:     m.session.Copy(),
		Unordered:   m.Unordered,
		PoolTimeout: m.PoolTimeout,
		Mode:        m.Mode,
	}
}

//SetUnordered bulk 錯誤不中斷
func (m *MgoInfo) SetUnordered() {

	m.Unordered = true
}

//SetMode  Mode changes the consistency mode for the session
func (m *MgoInfo) SetMode(mode mgo.Mode) {

	m.session.SetMode(mode, true)
	m.Mode = mode
}
