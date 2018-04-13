package goblog

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
)

func Install() error {
	sqlserver := fmt.Sprintf("%s:%s@tcp(%s:%s)/", Conf.Db.Username, Conf.Db.Password, Conf.Db.Address, Conf.Db.Port)
	var err error
	DB, err = sql.Open(Conf.Db.Driver, sqlserver)
	if err != nil {
		return err
	}

	DB.Exec("CREATE DATABASE " + Conf.Db.Dbname)

	DB.Exec("use " + Conf.Db.Dbname)

	_, err = DB.Exec(`CREATE TABLE article(
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20) NOT NULL,
		uuid VARCHAR(40) NOT NULL,
		tag CHAR(10) NOT NULL,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		edit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)
	if err != nil {
		return err
	}

	if _, err = os.Stat("article"); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll("article", 0775)
		}
	}

	if err != nil {
		return err
	}

	err = AddArticle("世界，您好！", "test", "欢迎使用goblog。这是您的第一篇文章。编辑或删除它，然后开始写作吧！")
	if err != nil {
		return err
	}

	passwd := fmt.Sprintf("%x", md5.Sum([]byte(Conf.Password)))
	passwd = fmt.Sprintf("%x", md5.Sum([]byte(passwd)))
	Conf.Password = passwd

	err = WriteConf("config.json")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("goblog.lock", []byte("lock"), 0755)
}
