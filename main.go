package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
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
	port    = flag.Int("p", 8080, "port")
	install = flag.Bool("i", false, "install goblog")
	help    = flag.Bool("h", false, "help")
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("goblog"))
	router.Use(sessions.Sessions("goblog-session", store))

	apiRouter := router.Group("/api")

	{
		adminC := &controller.Admin{}
		apiRouter.POST("/login", adminC.Login)
		apiRouter.POST("/logout", adminC.Logout)
		apiRouter.GET("/info", adminC.GetInfo)
		apiRouter.PATCH("/info", loginMiddleware, adminC.PatchInfo)
		apiRouter.GET("/logo", adminC.GetLogo)
		apiRouter.PUT("/logo", loginMiddleware, adminC.PutLogo)
		apiRouter.GET("/assets/:uuid", adminC.GetAssets)
		apiRouter.PUT("/assets", loginMiddleware, adminC.PutAssets)
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
	if *install {
		log.Println(installGoBlog())
		os.Exit(0)
	}

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

	if _, err := os.Stat(core.Conf.DataDir + "/article"); os.IsNotExist(err) {
		os.MkdirAll(core.Conf.DataDir+"/article", 0755)
	}

	if _, err := os.Stat(core.Conf.DataDir + "/assets"); os.IsNotExist(err) {
		os.MkdirAll(core.Conf.DataDir+"/assets", 0755)
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
		ctx.JSON(200, core.Result(core.ResultUnauthorizedCode, "unauthorized"))
		ctx.Abort()
	}
}

func installGoBlog() error {
	conf := core.Config{}
	fmt.Println("======================")
	fmt.Println("首次运行需要部署一些东西")
	fmt.Println("======================")
	fmt.Print("请输入博客名:")
	fmt.Scanln(&conf.Name)
	fmt.Print("请输入登录密码:")
	fmt.Scanln(&conf.Password)
	md5Data := md5.Sum([]byte(conf.Password))
	sha1Data := sha1.Sum([]byte(md5Data[:]))
	conf.Password = hex.EncodeToString(sha1Data[:])
	fmt.Println("是否开启分类功能(y/n)")
	tmp := "n"
	fmt.Scanln(&tmp)
	if tmp == "y" {
		conf.UseCategory = true
	}
	fmt.Println("======================")
	fmt.Println("数据库设置")
	fmt.Println("======================")
	fmt.Println("选择数据库驱动")
	fmt.Println("1.mysql")
	fmt.Println("2.postgreSQL")
	fmt.Scanln(&tmp)
	if tmp == "2" {
		conf.Db.Driver = "postgres"
	} else {
		conf.Db.Driver = "mysql"
	}
	fmt.Println("数据库用户名(root)")
	fmt.Scanln(&conf.Db.Username)
	if conf.Db.Username == "" {
		conf.Db.Username = "root"
	}
	fmt.Println("数据库密码")
	fmt.Scanln(&conf.Db.Password)
	fmt.Println("数据库地址(localhost)")
	fmt.Scanln(&conf.Db.Address)
	if conf.Db.Address == "" {
		conf.Db.Address = "localhost"
	}
	fmt.Println("数据库端口(3306)")
	fmt.Scanln(&conf.Db.Port)
	if conf.Db.Port == "" {
		conf.Db.Port = "3306"
	}
	fmt.Println("数据库名(goblog)")
	fmt.Scanln(&conf.Db.Dbname)
	if conf.Db.Dbname == "" {
		conf.Db.Dbname = "goblog"
	}

	conf.DataDir = "data"
	conf.Mode = "release"
	core.Conf = &conf
	return core.SaveConf()
}
