package goblog

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

func ParseConf(config Config) (Config,error) {

}