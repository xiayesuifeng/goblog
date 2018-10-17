package main

import (
	"flag"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/article"
	"gitlab.com/xiayesuifeng/goblog/category"
	"gitlab.com/xiayesuifeng/goblog/controller"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/database"
	"log"
	"os"
	"strconv"
)

var (
	port = flag.Int("p", 8080, "port")
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("goblog"))
	router.Use(sessions.Sessions("goblog-session", store))

	apiRouter := router.Group("/api")

	apiRouter.GET("/info", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"name":        core.Conf.Name,
			"useCategory": core.Conf.UseCategory,
		})
	})

	{
		adminC := &controller.Admin{}
		apiRouter.POST("/login", adminC.Login)
		apiRouter.POST("/logout", adminC.Logout)
	}

	{
		categoryC := &controller.Category{}
		apiRouter.GET("/category", categoryC.Gets)
		apiRouter.GET("/category/:id", categoryC.Get)
		apiRouter.POST("/category", loginMiddleware, categoryC.Post)
		apiRouter.PUT("/category/:id", loginMiddleware, categoryC.Put)
		apiRouter.DELETE("/category/:id", loginMiddleware, categoryC.Delete)
	}

	{
		tagC := &controller.Tag{}
		apiRouter.GET("/tag", tagC.Gets)
		apiRouter.GET("/tag/:tag", tagC.Get)
	}

	{
		articleC := &controller.Article{}
		apiRouter.GET("/article", articleC.Gets)
		apiRouter.GET("/article/id/:id", articleC.Get)
		apiRouter.GET("/article/category/:category_id", articleC.GetByCategory)
		apiRouter.GET("/article/uuid/:uuid/:mode", articleC.GetByUuid)
		apiRouter.POST("/article", loginMiddleware, articleC.Post)
		apiRouter.PUT("/article/:id", loginMiddleware, articleC.Put)
		apiRouter.DELETE("/article/:id", loginMiddleware, articleC.Delete)
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

	if _, err := os.Stat(core.Conf.DataDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(core.Conf.DataDir, 0755)
		} else {
			log.Panicln("data dir create failure")
		}
	}

	gin.SetMode(core.Conf.Mode)

	if !core.Conf.UseCategory {
		tmp := category.Category{Name: "other"}
		if db.Where(&tmp).First(&tmp).RecordNotFound() {
			if err := db.Create(&tmp).Error; err != nil {
				log.Panicln(err)
			}
		}

		db.Model(&article.Article{}).Updates(article.Article{CategoryId: tmp.ID})
		core.Conf.OtherCategoryId = tmp.ID
	}
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
