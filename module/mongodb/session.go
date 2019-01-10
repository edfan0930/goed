package mongodb

import (
	"runtime"

	"gopkg.in/mgo.v2"
)

//MgoInfo Control session
type (
	mgoInfo struct {
		session *mgo.Session
	}
)

// NewMgoInfo --
// 當GC啟動時,執行finalizer
func NewMgoInfo() (m *mgoInfo) {
	m = &mgoInfo{
		session: session.Copy(),
	}
	runtime.SetFinalizer(m, finalizer)
	return
}

//Close Session
func finalizer(m *mgoInfo) {
	m.session.Close()
}
