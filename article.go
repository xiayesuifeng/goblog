package goblog

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
)

type Article struct {
	Name string
	Uuid string
	Tag string
	CreateTime int64
	EditTime int64
}

func AddArticle(name string,tag string,context string) error {
	md_uuid := uuid.NewV1().String()
	_,err := DB.Exec(`INSERT INTO article (name,uuid,tag) VALUES (?,?,?)`,name,md_uuid,tag)
	if err!= nil {
		return err
	}

	return ioutil.WriteFile("article/"+md_uuid+".md",[]byte(context),0644)
}

func UpdateArticle(name string,tag string,context string) error{
	return nil
}

func DelArticle(name string) error{
	var md_uuid string
	row:=DB.QueryRow("SELECT uuid FROM article WHERE name=?",name)
	err := row.Scan(&md_uuid)
	if err!= nil {
		return err
	}

	if err=os.Remove("article/"+md_uuid+".md"); err!= nil{
		return err
	}
	_,err = DB.Exec("DELETE FROM article WHERE name=?",name)
	return err
}