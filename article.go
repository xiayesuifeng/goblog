package goblog

import (
	"fmt"
	"database/sql"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

func AddArticle(name string,tag string,context string,config Config) error {
	sqlserver := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",config.Db.Username,config.Db.Password,config.Db.Address,config.Db.Port,config.Db.Dbname)
	db,err := sql.Open(config.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}

	md_uuid := uuid.NewV1().String()
	_,err = db.Exec(`INSERT INTO article (name,uuid,tag) VALUES (?,?,?)`,name,md_uuid,tag)
	if err!= nil {
		return err
	}

	return ioutil.WriteFile("article/"+md_uuid+".md",[]byte(context),0644)
}

func UpdateArticle(name string,tag string,context string,config Config) error{
	return nil
}

func DelArticle(name string,config Config) error{
	return nil
}