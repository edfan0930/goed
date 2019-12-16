package mongodb

import (
	"github.com/globalsign/mgo/bson"
)

const (
	//DBUser db name
	DBUser = "user"
	//CollectionUser collection name
	CollectionUser = "file"
)

//Field const
const (
	FieldID      = "id"
	FieldName    = "name"
	FieldAccount = "account"
	FieldAge     = "age"
)

//User filed struct
type User struct {
	ID      bson.ObjectId
	Account string
	Name    string
	Age     uint
}

//NewUser instance user
func NewUser() *User {
	return &User{}
}

//Rename ...
func (u *User) Rename(m *MgoInfo, id bson.ObjectId, name string) error {

	find := bson.M{
		FieldID: id,
	}
	update := bson.M{
		"$set": bson.M{
			FieldName: name,
		},
	}

	return m.update(DBUser, CollectionUser, find, update)
}
