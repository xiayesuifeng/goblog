package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() error {
	db := Conf.Db
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		db.Username, db.Password, db.Address, db.Port, db.Dbname)
	var err error
	DB,err = gorm.Open("mysql", args)
	return err
}
