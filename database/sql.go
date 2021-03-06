package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"log"
	"sync"
)

var db *gorm.DB
var once sync.Once

func initDB() {
	config := conf.Conf.Db
	args := ""

	switch config.Driver {
	case "mysql":
		args = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
			config.Username, config.Password, config.Address, config.Port, config.Dbname)
	case "postgres":
		args = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
			config.Address, config.Port, config.Dbname, config.Username, config.Password)
	case "sqlite3":
		args = fmt.Sprintf("%s/%s.db?_auth&_auth_user=%s&_auth_pass=%s", conf.Conf.DataDir, config.Dbname, config.Username, config.Password)
	case "mssql":
		args = fmt.Sprintf("sqlserver://%s:%s@localhost:%s?database=%s", config.Username, config.Password, config.Port, config.Dbname)
	}

	var err error
	db, err = gorm.Open(config.Driver, args)
	if err != nil {
		log.Fatalln(err)
	}
}

func Instance() *gorm.DB {
	once.Do(initDB)
	return db
}
