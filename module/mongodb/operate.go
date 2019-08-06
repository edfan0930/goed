package mongodb

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	//SizeNotEqual  size not equal condition of DB find bulk
	SizeNotEqual = "size not equal"
	//ModifiedLess  modified less than condition of DB update
	ModifiedLess = "modified less"
	//NotFound data not found
	NotFound = "not found"
)

type (
	//SQL DB syntax
	sql struct {
		Selector bson.M
		Update   bson.M
	}
)

//newSQL return instance
func newSQL(update bson.M) *sql {
	return &sql{
		Update: update,
	}
}

//insertByBulk
func (m *MgoInfo) insertByBulk(db, collection string, i []interface{}) (*mgo.BulkResult, error) {

	bulk := m.session.DB(db).C(collection).Bulk()

	//設定錯誤時的處理
	if m.Unordered == true {
		bulk.Unordered()
	}

	bulk.Insert(i...)

	return bulk.Run()
}

//updateByBulk
func (m *MgoInfo) updateByBulk(db, collection string, mgosql []*sql) (*mgo.BulkResult, error) {

	bulk := m.session.DB(db).C(collection).Bulk()

	if m.Unordered == true {
		bulk.Unordered()
	}

	for k := range mgosql {
		bulk.Update(mgosql[k].Selector, mgosql[k].Update)
	}

	return bulk.Run()
}

//insert ...
func (m *MgoInfo) insert(db, collection string, i interface{}) error {

	return m.session.DB(db).C(collection).Insert(i)
}

//update ...
func (m *MgoInfo) update(db, collection string, condition, update bson.M) (err error) {

	c := m.session.DB(db).C(collection)

	_, err = c.UpdateAll(condition, update)

	return
}

//removeAll ...
func (m *MgoInfo) removeAll(db, collection string, condition bson.M) (err error) {

	c := m.session.DB(db).C(collection)

	_, err = c.RemoveAll(condition)

	return err
}

//findOne ...
func (m *MgoInfo) findOne(db, collection string, condition bson.M, i interface{}) error {

	c := m.session.DB(db).C(collection)

	return c.Find(condition).One(&i)
}

//findAll find all
func (m *MgoInfo) findAll(db, collection string, condition bson.M, i interface{}) error {

	c := m.session.DB(db).C(collection)

	return c.Find(condition).All(i)
}

//Ping ping session
func (m *MgoInfo) Ping() error {

	return m.session.Ping()
}
