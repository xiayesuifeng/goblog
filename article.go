package goblog

import (
	"fmt"
	"database/sql"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
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
	sqlserver := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",config.Db.Username,config.Db.Password,config.Db.Address,config.Db.Port,config.Db.Dbname)
	db,err := sql.Open(config.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}

	var md_uuid string
	row:=db.QueryRow("SELECT uuid FROM article WHERE name=?",name)
	err = row.Scan(&md_uuid)
	if err!= nil {
		return err
	}

	if err=os.Remove("article/"+md_uuid+".md"); err!= nil{
		return err
	}
	_,err = db.Exec("DELETE FROM article WHERE name=?",name)
	return err
}