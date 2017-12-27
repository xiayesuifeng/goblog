package goblog

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"time"
	"errors"
)

type Article struct {
	Id string
	Name string
	Uuid string
	Tag string
	CreateTime int64
	EditTime int64
}

func AddArticle(name string,tag string,context string) error {
	row:=DB.QueryRow("SELECT id FROM article WHERE name=?",name)
	var id int
	err:=row.Scan(&id)
	if err==nil {
		return errors.New("name exist")
	}

	md_uuid := uuid.NewV1().String()
	_,err = DB.Exec(`INSERT INTO article (name,uuid,tag) VALUES (?,?,?)`,name,md_uuid,tag)
	if err!= nil {
		return err
	}

	return ioutil.WriteFile("article/"+md_uuid+".md",[]byte(context),0644)
}

func UpdateArticle(oldName,newName,tag,context string) error{
	row:=DB.QueryRow("SELECT uuid FROM article WHERE name=?",oldName)
	var u string
	row.Scan(&u)

	out,err := os.OpenFile("article/"+u+".md",os.O_RDWR | os.O_TRUNC,0)
	if err!= nil {
		return err
	}
	defer out.Close()
	_,err=out.WriteString(context)
	if err != nil {
		return err
	}

	lastTime := time.Now().Format("2006-01-02 15:04:05")
	_,err=DB.Exec("UPDATE article SET name=?,tag=?,edit_time=? WHERE name=?",newName,tag,lastTime,oldName)
	return err
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