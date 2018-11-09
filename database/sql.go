package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/xiayesuifeng/goblog/core"
	"log"
	"sync"
)

var db *gorm.DB
var once sync.Once

func initDB() {
	conf := core.Conf.Db
	args := ""

	switch conf.Driver {
	case "mysql":
		args = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
			conf.Username, conf.Password, conf.Address, conf.Port, conf.Dbname)
	case "postgres":
		args = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			conf.Address, conf.Port, conf.Dbname, conf.Username, conf.Password)
	}

	var err error
	db, err = gorm.Open(conf.Driver, args)
	if err != nil {
		log.Fatalln(err)
	}
}

func Instance() *gorm.DB {
	once.Do(initDB)
	return db
}
