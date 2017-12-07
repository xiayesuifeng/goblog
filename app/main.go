package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/1377195627/goblog"
	"log"
	"fmt"
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

	api.GET("/name/:name", func(context *gin.Context) {

	})

	api.GET("/tag/:tag", func(context *gin.Context) {

	})

	route.GET("/install", func(context *gin.Context) {
		context.HTML(http.StatusOK,"install.html",gin.H{})
	})

	route.POST("/install", func(context *gin.Context) {
		err = context.Bind(&config)
		if err != nil {
			log.Fatal(err)
		}

		context.JSON(http.StatusOK,config)
	})

	route.Run()
}
