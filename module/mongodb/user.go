package mongodb

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	user struct {
		Name   string
		Gender string
	}
)

//CollectionUser collection name
const CollectionUser = "user"

//NewUser return Instance
func NewUser(name, gender string) *user {
	return &user{
		Name:   name,
		Gender: gender,
	}
}

//NewUserInterface return Interface Instance
func NewUserInterface(u ...*user) []interface{} {
	ui := make([]interface{}, len(u))
	for k := range u {
		fmt.Println("users is", *u[k])
		ui = append(ui, *u[k])
	}
	return ui
}

//UsersInsertByBulk 批次處理 insert
func (m *mgoInfo) UsersInsertByBulk(u ...user) (err error) {
	bulk := m.session.DB("").C(CollectionUser).Bulk()
	bulk.Unordered()
	for k := range u {
		bulk.Insert(u[k])
	}
	_, err = bulk.Run()
	if err != nil {
		log.Println(err)
	}
	return
}

//UserUpdateByBulk 批次處理 update
func (m *mgoInfo) UserUpdateByBulk(u ...user) (*mgo.BulkResult, error) {
	bulk := m.session.DB("").C(CollectionUser).Bulk()
	bulk.Unordered()
	for k := range u {
		bulk.Update(
			bson.M{"name": u[k].Name, "gender": "male"},
			bson.M{"$set": bson.M{"gender": u[k].Gender}},
		)
	}
	result, bulkErr := bulk.Run()
	if bulkErr != nil {
		log.Println(bulkErr)
	}
	return result, bulkErr
}
