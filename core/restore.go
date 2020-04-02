package core

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func Restore(file string, useOldConfig bool) error {
	config, err := GetConfigForZip(file)
	if err != nil {
		return err
	}

	if !useOldConfig {
		config.DataDir = Conf.DataDir
		config.Db = Conf.Db
	}

	Conf = config
	if err := SaveConf(); err != nil {
		return err
	}

	if err := Unzip(file, Conf.DataDir); err != nil {
		return err
	}

	log.Println("Restore Database")
	return RestoreDatabase()
}

func RestoreDatabase() error {
	if Conf.Db.Driver == "mysql" {
		if err := RestoreMysqlDatabase(); err != nil {
			return err
		}
	} else if Conf.Db.Driver == "postgres" {
		if err := RestorePostgresDatabase(); err != nil {
			return err
		}
	}

	os.Remove(Conf.DataDir + "/database.sql")

	return nil
}

func RestoreMysqlDatabase() error {
	path, err := exec.LookPath("mysql")
	if err != nil {
		return err
	}

	sql, err := os.Open(Conf.DataDir + "/database.sql")
	if err != nil {
		return err
	}
	defer sql.Close()

	db := Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-P", db.Port, "-u", db.Username, "-p"+db.Password, db.Dbname)
	cmd.Stdin = sql

	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return nil
}

func RestorePostgresDatabase() error {
	path, err := exec.LookPath("pg_restore")
	if err != nil {
		return err
	}

	db := Conf.Db

	cmd := exec.Command(path, "-h", db.Address, "-p", db.Port, "-U", db.Username, "-d", db.Dbname, Conf.DataDir+"/backup/database.sql")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PGPASSWORD="+db.Password)

	data, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(data))
	}

	return nil
}
