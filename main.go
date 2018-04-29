package main

import (
	"flag"
	"fmt"
	"github.com/1377195627/goblog/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	once sync.Once
	port = flag.Int("p", 8080, "port")
)

func main() {
	flag.Parse()

	err := core.ParseConf("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("请配置config.json")
			os.Exit(0)
		}
		log.Panicln(err)
	}

	if err := core.InitDB();err!=nil{
		log.Panicln(err)
	}

	router := gin.Default()

	store := sessions.NewCookieStore([]byte("goblog"))
	router.Use(sessions.Sessions("goblog-session", store))

	router.Run(":" + strconv.Itoa(*port))
}
