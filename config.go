package goblog

import (
	"encoding/json"
	"os"
)

var Conf *Config

type Config struct {
	Name string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
	Db Database `json:"database"`
}

type Database struct {
	Driver string `json:"driver" form:"driver"`
	Address string `json:"address" form:"address"`
	Port string `json:"port" form:"port"`
	Dbname string `json:"dbname" form:"dbname"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"db_password"`
}

func ParseConf(config string) error {
	var c Config

	conf,err := os.Open(config)
	if err!= nil {
		return err
	}
	err = json.NewDecoder(conf).Decode(&c)

	Conf = &c
	return err
}

func WriteConf(outPath string) error{
	out,err := os.OpenFile(outPath,os.O_RDWR | os.O_TRUNC,0)
	if err!= nil {
		return err
	}
	return json.NewEncoder(out).Encode(Conf)
}