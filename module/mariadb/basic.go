package mariadb

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	//DIALECT
	DIALECT = "mysql"
)

//Maria ...
type Maria struct {
	db *gorm.DB
}

//NewMaria instance Maria struct
//initialize a new db connection
//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
//host ip:port
func NewMaria(user, password, host, dbName string) (*Maria, error) {

	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=UTC&timeout=30s", user, password, host, dbName)

	//"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=UTC&timeout=30s"
	db, err := gorm.Open(DIALECT, args)

	return &Maria{
		db: db,
	}, err
}

//Close close current db connection.
func (m *Maria) Close() {
	m.db.Close()
}
