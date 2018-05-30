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
	"github.com/1377195627/goblog/db"
)

var (
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

	db.Init()

	router := gin.Default()

	store := sessions.NewCookieStore([]byte("goblog"))
	router.Use(sessions.Sessions("goblog-session", store))

	apiRouter := router.Group("/api")

	apiRouter.POST("/login")

	apiRouter.GET("/tag")
	apiRouter.GET("/tag/:tag")

	apiRouter.GET("/article/name/:name")
	apiRouter.GET("/article/uuid/:uuid")
	apiRouter.POST("/article")
	apiRouter.PUT("/article")
	apiRouter.DELETE("/article/:name")

	router.Run(":" + strconv.Itoa(*port))
}
