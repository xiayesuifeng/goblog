package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.com/xiayesuifeng/goblog/core"
	"log"
	"sync"
)

var db *gorm.DB
var once sync.Once

func initDB() {
	conf := core.Conf.Db
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		conf.Username, conf.Password, conf.Address, conf.Port, conf.Dbname)
	var err error
	db, err = gorm.Open("mysql", args)
	if err != nil {
		log.Fatalln(err)
	}
}

func Instance() *gorm.DB {
	once.Do(initDB)
	return db
}
