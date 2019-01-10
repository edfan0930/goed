package mongodb

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

// Database -- Database Name
const Database = "test"

//var mongoDialInfo *mgo.DialInfo

var session *mgo.Session

var hostList, userName, password, database string

var host []string

func init() {
	SetDialInfo()
}

// SetDialInfo init mgo session
func SetDialInfo() {
	var err error
	host = []string{"127.0.0.1"}
	timeOut := 30
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    host,
		Timeout:  time.Duration(timeOut) * time.Second,
		Database: Database,
		Username: userName,
		Password: password,
	}
	session, err = mgo.DialWithInfo(mongoDialInfo)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	//slave 為主
	//session.SetMode(mgo.Eventual, true)
}
