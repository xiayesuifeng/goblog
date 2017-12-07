package goblog

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name string `json:"name" form:"name"`
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

func ParseConf(config string) (Config,error) {
	var c Config

	conf,err := os.Open(config)
	if err!= nil {
		return c,err
	}
	err = json.NewDecoder(conf).Decode(&c)

	return c,err
}

func WriteConf(config Config,outPath string) error{
	out,err := os.OpenFile(outPath,os.O_RDWR | os.O_TRUNC,0)
	if err!= nil {
		return err
	}
	return json.NewEncoder(out).Encode(config)
}