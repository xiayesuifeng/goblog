package core

import (
	"encoding/json"
	"os"
)

var Conf *Config

type Config struct {
	Mode string `json:"mode"`
	Name string `json:"name"`
	Password string `json:"password"`
	UseCategory bool `json:"useCategory"`
	DataDir string `json:"data_dir"`
	Db   Database `json:"database"`
	Smtp Smtp     `json:"smtp"`
}

type Database struct {
	Address  string `json:"address" form:"address"`
	Port     string `json:"port" form:"port"`
	Dbname   string `json:"dbname" form:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Smtp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

func ParseConf(config string) error {
	var c Config

	conf, err := os.Open(config)
	if err != nil {
		return err
	}
	err = json.NewDecoder(conf).Decode(&c)

	Conf = &c
	return err
}
