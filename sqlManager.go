package goblog

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var DB *sql.DB

func InitSql(config Config) error{
	sqlserver := fmt.Sprintf("%s:%s@tcp(%s:%s)/",config.Db.Username,config.Db.Password,config.Db.Address,config.Db.Port)
	var err error
	DB,err = sql.Open(config.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}
	return nil
}