package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/1377195627/goblog"
	"log"
	"github.com/1377195627/goblog/api"
	"flag"
	"strconv"
)

var port = flag.Int("p",80,"listen port,default 8080")
var server = flag.String("S","127.0.0.1","listen server,default 127.0.0.1")

func main() {
	Init()

	router := gin.Default()
	//router.LoadHTMLGlob("view/*")
	router.Static("/images", "static/images")

	apiRoter := router.Group("api")
	apiRoter.POST("/install",goblog.InstallRouter)
	apiRoter.GET("/name",api.Name)
	apiRoter.POST("/login", api.Login)
	apiRoter.GET("/tags", api.Tags)
	apiRoter.GET("/tag", api.Tag)
	apiRoter.GET("/tag/:tag", api.TagBytag)
	apiRoter.GET("/article/name/:name", api.ArticleByName)
	apiRoter.GET("/article/uuid/:uuid", api.ArticleByUuid)
	apiRoter.POST("/article/new", api.ArticleNew)
	apiRoter.PUT("/article/edit", api.ArticleEdit)
	apiRoter.DELETE("/article/del/:name", api.ArticleDel)

	//router.GET("/install", goblog.InstallRouter)
	//router.POST("/install", goblog.InstallRouter)
	router.GET("/", goblog.HomeRouter)
	router.Run(*server+":"+strconv.Itoa(*port))
}

func Init() {
	err := goblog.ParseConf("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = goblog.InitSql()
	if err != nil {
		log.Fatal(err)
	}
}
