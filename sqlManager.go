package goblog

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
)

var DB *sql.DB

func InitSql() error{
	conf := Conf
	var sqlserver string
	if _, err := os.Stat("goblog.lock"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	sqlserver = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",conf.Db.Username,conf.Db.Password,conf.Db.Address,conf.Db.Port,conf.Db.Dbname)
	var err error
	DB,err = sql.Open(conf.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}
	return nil
}