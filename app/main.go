package main

import "gopkg.in/gin-gonic/gin.v1"

func main() {
	route := gin.Default()

	api := route.Group("api")

	api.GET("/:name", func(context *gin.Context) {

	})

	api.GET("/tag/:tag", func(context *gin.Context) {

	})

	route.POST("/install", func(context *gin.Context) {

	})
}
