package goblog

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func Install(config Config) error {
	sqlserver := fmt.Sprintf("%s:%s@%s:%s/",config.Db.Username,config.Db.Password,config.Db.Address,config.Db.Port)
	db,err := sql.Open(config.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}

	db.Exec("CREATE DATABASE "+config.Db.Dbname+" IF NOT EXIST")

	db.Exec("use "+config.Db.Dbname)

	_,err=db.Exec(`CREATE TABLE article(
		name VARCHAR(20) NOT NULL
		md_name VARCHAR(20) NOT NULL
		tag CHAR(10) NOT NULL
		create_time DATETIME NOT NULL
		edit_time DATEIME NOT NULL
)`)

	if err!= nil {
		return err
	}

	return nil
}
