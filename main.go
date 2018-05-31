package main

import (
	"flag"
	"github.com/1377195627/goblog/article"
	"github.com/1377195627/goblog/category"
	"github.com/1377195627/goblog/controller"
	"github.com/1377195627/goblog/core"
	"github.com/1377195627/goblog/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

var (
	port = flag.Int("p", 8080, "port")
)

func main() {
	router := gin.Default()

	store := sessions.NewCookieStore([]byte("goblog"))
	router.Use(sessions.Sessions("goblog-session", store))

	apiRouter := router.Group("/api")

	apiRouter.POST("/login")

	{
		category := &controller.Category{}
		apiRouter.GET("/category", category.Gets)
		apiRouter.GET("/category/:id", category.Get)
		apiRouter.POST("/category", category.Post)
		apiRouter.PUT("/category/:id", category.Put)
		apiRouter.DELETE("/category/:id", category.Delete)
	}

	apiRouter.GET("/tag")
	apiRouter.GET("/tag/:tag")

	{
		article := &controller.Article{}
		apiRouter.GET("/article", article.Gets)
		apiRouter.GET("/article/category/:category", article.GetByCategory)
		apiRouter.GET("/article/name/:name", article.Get)
		apiRouter.GET("/article/uuid/:uuid", article.GetByUuid)
		apiRouter.POST("/article", article.Post)
		apiRouter.PUT("/article", article.Put)
		apiRouter.DELETE("/article/:id", article.Delete)
	}

	router.Run(":" + strconv.Itoa(*port))
}

func init() {
	flag.Parse()

	err := core.ParseConf("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("please config config.json")
			os.Exit(0)
		}
		log.Panicln(err)
	}

	db := database.Instance()
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&article.Article{})
}
