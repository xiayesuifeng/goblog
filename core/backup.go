package core

import (
	"errors"
	"fmt"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func Backup() error {
	fmt.Println("温馨提示：使用备份功能为保证数据完整性，请停止运行 GoBlog 后再进行，请在确认后输入Y继续(Y/N):")
	input := ""
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}

	if input != "Y" && input != "y" {
		return nil
	}

	if _, err := os.Stat(conf.Conf.DataDir + "/backup"); os.IsNotExist(err) {
		os.MkdirAll(conf.Conf.DataDir+"/backup", 0755)
	}

	fmt.Println("Dumping database")
	if err := DumpDatabase(); err != nil {
		return err
	}

	zipName := conf.Conf.DataDir + time.Now().Format("/backup/Backup-GoBlog-20060102150405.zip")
	if err := Zip(conf.Conf.DataDir, zipName); err != nil {
		return err
	}

	os.Remove(conf.Conf.DataDir + "/backup/database.sql")

	fmt.Println("Backup save to ", zipName)

	return nil
}

func DumpDatabase() error {
	if conf.Conf.Db.Driver == "mysql" {
		err := DumpMysqlDatabase()
		if err != nil {
			return err
		}
	} else if conf.Conf.Db.Driver == "postgres" {
		err := DumpPostgresDatabase()
		if err != nil {
			return err
		}
	}

	return nil
}

func DumpMysqlDatabase() error {
	path, err := exec.LookPath("mysqldump")
	if err != nil {
		return err
	}

	db := conf.Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-P", db.Port, "-u", db.Username, "-p"+db.Password, db.Dbname)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return ioutil.WriteFile(conf.Conf.DataDir+"/backup/database.sql", data, 0664)
}

func DumpPostgresDatabase() error {
	path, err := exec.LookPath("pg_dump")
	if err != nil {
		return err
	}

	db := conf.Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-p", db.Port, "-U", db.Username, "-f", conf.Conf.DataDir+"/backup/database.sql", db.Dbname)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+db.Password)

	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return nil
}
