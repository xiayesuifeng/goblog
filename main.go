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

	{
		adminC := &controller.Admin{}
		apiRouter.POST("/login", adminC.Login)
		apiRouter.POST("/logout", adminC.Logout)
	}

	{
		categoryC := &controller.Category{}
		apiRouter.GET("/category", categoryC.Gets)
		apiRouter.GET("/category/:id", categoryC.Get)
		apiRouter.POST("/category",loginMiddleware, categoryC.Post)
		apiRouter.PUT("/category/:id",loginMiddleware, categoryC.Put)
		apiRouter.DELETE("/category/:id",loginMiddleware, categoryC.Delete)
	}

	{
		tagC := &controller.Tag{}
		apiRouter.GET("/tag",tagC.Gets)
		apiRouter.GET("/tag/:tag",tagC.Get)
	}

	{
		articleC := &controller.Article{}
		apiRouter.GET("/article", articleC.Gets)
		apiRouter.GET("/article/category/:category", articleC.GetByCategory)
		apiRouter.GET("/article/name/:name", articleC.Get)
		apiRouter.GET("/article/uuid/:uuid/:mode", articleC.GetByUuid)
		apiRouter.POST("/article",loginMiddleware, articleC.Post)
		apiRouter.PUT("/article/:id",loginMiddleware, articleC.Put)
		apiRouter.DELETE("/article/:id",loginMiddleware, articleC.Delete)
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

func loginMiddleware(ctx *gin.Context) {
	session := sessions.Default(ctx)

	login := session.Get("login")
	if login == nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "unauthorized",
		})
		ctx.Abort()
	}
}