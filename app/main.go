package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/1377195627/goblog"
	"log"
)

var config goblog.Config

func main() {
	var err error
	config,err = goblog.ParseConf("config.json")
	if err != nil {
		log.Fatal(err)
	}

	route := gin.Default()
	route.LoadHTMLGlob("view/*")

	api := route.Group("api")

	api.GET("/tag/:tag", func(context *gin.Context) {

	})

	api.GET("/article/:name", func(context *gin.Context) {

	})

	api.POST("/article/new", func(context *gin.Context) {

	})

	api.POST("/article/del", func(context *gin.Context) {

	})

	api.PUT("/article/edit", func(context *gin.Context) {

	})

	route.GET("/install", func(context *gin.Context) {
		context.HTML(http.StatusOK,"install.html",gin.H{})
	})

	route.POST("/install", func(context *gin.Context) {
		err = context.Bind(&config)
		if err != nil {
			log.Fatal(err)
		}

		if err=goblog.Install(config);err!=nil{
			log.Fatal(err)
			context.String(http.StatusOK,"安装失败" )
		}

		context.String(http.StatusOK,"安装完成" )

	})

	route.Run()
}
