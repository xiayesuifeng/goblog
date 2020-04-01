package core

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func DumpDatabase() error {
	if Conf.Db.Driver == "mysql" {
		err := DumpMysqlDatabase()
		if err != nil {
			return err
		}
	} else if Conf.Db.Driver == "postgres" {
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

	db := Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-P", db.Port, "-u", db.Username, "-p"+db.Password, db.Dbname)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return ioutil.WriteFile(Conf.DataDir+"/backup/database.sql", data, 0664)
}

func DumpPostgresDatabase() error {
	path, err := exec.LookPath("pg_dump")
	if err != nil {
		return err
	}

	db := Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-p", db.Port, "-U", db.Username, "-f", Conf.DataDir+"/backup/database.sql", db.Dbname)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+db.Password)

	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return nil
}
