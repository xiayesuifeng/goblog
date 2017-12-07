package goblog

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name string `json:"name"`
	Db Database `json:"database"`
}

type Database struct {
	Driver string `json:"driver"`
	Address string `json:"address"`
	Port string `json:"port"`
	Dbname string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
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