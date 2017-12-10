package goblog

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"crypto/md5"
	"os"
	"io/ioutil"
)

func Install(config Config) error {
	sqlserver := fmt.Sprintf("%s:%s@tcp(%s:%s)/",config.Db.Username,config.Db.Password,config.Db.Address,config.Db.Port)
	db,err := sql.Open(config.Db.Driver,sqlserver)
	if err!= nil {
		return err
	}

	fmt.Println(db.Exec("CREATE DATABASE "+config.Db.Dbname))

	db.Exec("use "+config.Db.Dbname)

	_,err=db.Exec(`CREATE TABLE article(
		name VARCHAR(20) NOT NULL,
		uuid VARCHAR(40) NOT NULL,
		tag CHAR(10) NOT NULL,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		edit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)

	if _,err = os.Stat("article"); err != nil {
		if os.IsNotExist(err) {
			err=os.MkdirAll("article", 0775)
		}
	}

	if err!= nil {
		return err
	}

	err = AddArticle("世界，您好！","test","欢迎使用goblog。这是您的第一篇文章。编辑或删除它，然后开始写作吧！",config)
	if err!= nil {
		return err
	}

	passwd := fmt.Sprintf("%x",md5.Sum([]byte(config.Password)))
	passwd = fmt.Sprintf("%x",md5.Sum([]byte(passwd)))
	config.Password = passwd

	err = WriteConf(config,"config.json")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("goblog.lock",[]byte("lock"),0755)
}
