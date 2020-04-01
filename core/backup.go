package core

import (
	"errors"
	"io/ioutil"
	"os/exec"
)

func DumpMysqlDatabase() error {
	path, err := exec.LookPath("mysqldump")
	if err != nil {
		return err
	}

	db := Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-P", db.Port, "-u", db.Username, "-p"+db.Password, db.Dbname)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return ioutil.WriteFile(Conf.DataDir+"/backup/database.sql", data, 0664)
}
