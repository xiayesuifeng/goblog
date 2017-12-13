package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/1377195627/goblog"
	"log"
	"github.com/1377195627/goblog/api"
)

func main() {
	Init()

	route := gin.Default()
	route.LoadHTMLGlob("view/*")

	apiRoter := route.Group("api")
	apiRoter.POST("/login", api.Login)
	apiRoter.GET("/tags", api.Tags)
	apiRoter.GET("/tag", api.Tag)
	apiRoter.GET("/tag/:tag", api.TagBytag)
	apiRoter.GET("/article/name/:name", api.ArticleByName)
	apiRoter.GET("/article/uuid/:uuid", api.ArticleByUuid)
	apiRoter.POST("/article/new", api.ArticleNew)
	apiRoter.DELETE("/article/del", api.ArticleDel)
	apiRoter.PUT("/article/edit", api.ArticleEdit)

	route.GET("/install", goblog.InstallRouter)
	route.POST("/install", goblog.InstallRouter)
	route.GET("/", goblog.HomeRouter)
	route.Run()
}

func Init(){
	err := goblog.ParseConf("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = goblog.InitSql()
	if err != nil {
		log.Fatal(err)
	}
}
