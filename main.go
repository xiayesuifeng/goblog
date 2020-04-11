package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/cloudflare/tableflip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/article"
	"gitlab.com/xiayesuifeng/goblog/category"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"gitlab.com/xiayesuifeng/goblog/controller"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/database"
	"gitlab.com/xiayesuifeng/goblog/plugin"
	_ "gitlab.com/xiayesuifeng/goblog/sql-driver"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	backup      = flag.Bool("b", false, "backup goblog")
	restore     = flag.String("r", "", "restore goblog backup file")
	port        = flag.Int("p", 8080, "port")
	install     = flag.Bool("i", false, "install goblog")
	help        = flag.Bool("h", false, "help")
	config      = flag.String("c", "config.json", "Config file for goblog")
	pidFile     = flag.String("pid-file", "", "path to pid file")
	useHttpPort = flag.Bool("autotls-use-custom-http-port", false,
		"Allow http port to be specified with -p when autotls is enabled")
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

	{
		pluginC := &controller.Plugin{}
		apiRouter.GET("/plugin", pluginC.Gets)
		apiRouter.GET("/plugin/:name", loginMiddleware, pluginC.Download)
		apiRouter.DELETE("/plugin/:name", loginMiddleware, pluginC.Delete)
	}

	plugin.InitRouters(apiRouter)

	webPath := os.Getenv("GOBLOG_WEB_PATH")
	if webPath != "" {
		router.Use(static.Serve("/", static.LocalFile(webPath, false)))
		router.NoRoute(func(c *gin.Context) {
			if !strings.Contains(c.Request.RequestURI, "/api") {
				path := strings.Split(c.Request.URL.Path, "/")
				if len(path) > 1 {
					c.File(webPath + "/index.html")
					return
				}
			}
		})
	}

	upg, err := tableflip.New(tableflip.Options{
		PIDFile: *pidFile,
	})
	if err != nil {
		log.Panicln(err)
	}
	defer upg.Stop()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		for s := range sig {
			switch s {
			case syscall.SIGHUP:
				log.Println(upg.Upgrade())
			default:
				os.Remove(*pidFile)
				upg.Stop()
			}
		}
	}()

	address := ":" + strconv.Itoa(*port)
	if conf.Conf.Tls.Enable && !*useHttpPort {
		address = ":http"
	}
	ln, err := upg.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	httpServer := &http.Server{
		Handler: router,
	}

	go httpServer.Serve(ln)

	var httpsServer *http.Server
	if conf.Conf.Tls.Enable {
		httpsLn, err := upg.Listen("tcp", ":https")
		if err != nil {
			log.Fatalln(err)
		}
		defer httpsLn.Close()

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.Conf.Tls.Domain...),
			Cache:      autocert.DirCache(conf.Conf.DataDir + "/acme"),
		}

		httpsServer = &http.Server{
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
			Handler:   router,
		}

		go httpsServer.ServeTLS(httpsLn, "", "")
	}

	if err := upg.Ready(); err != nil {
		log.Panicln(err)
	}

	core.Upg = upg
	<-upg.Exit()

	gracefulShutdown(httpServer, httpsServer)
}

func gracefulShutdown(server ...*http.Server) {
	for _, s := range server {
		if s == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}
}

func init() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *install {
		log.Println(installGoBlog())
		os.Exit(0)
	}

	if *restore != "" {
		fmt.Println("温馨提示：使用恢复功能为保证数据完整性，请停止运行 GoBlog 后再进行，请在确认后输入Y继续(Y/N):")
		input := ""
		if _, err := fmt.Scanln(&input); err != nil {
			log.Println(err)
			os.Exit(0)
		}

		if input != "Y" && input != "y" {
			os.Exit(0)
		}

		fmt.Println("是否使用备份中的 DataDir 与 database 覆盖现有配置(Y/N):")
		if _, err := fmt.Scanln(&input); err != nil {
			log.Println(err)
			os.Exit(0)
		}

		if input == "Y" || input == "y" {
			if err := core.Restore(*restore, true); err != nil {
				log.Println(err)
			}
			log.Println("备份已成功恢复")
			os.Exit(0)
		}
	}

	err := conf.ParseConf(*config)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("please config config.json")
			os.Exit(0)
		}
		log.Panicln(err)
	}

	if _, err := os.Stat(conf.Conf.DataDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(conf.Conf.DataDir, 0755)
		} else {
			log.Panicln("data dir create failure")
		}
	}

	if *restore != "" {
		if err := core.Restore(*restore, false); err != nil {
			log.Println(err)
		}
		log.Println("备份已成功恢复")
		os.Exit(0)
	}

	if _, err := os.Stat(conf.Conf.DataDir + "/article"); os.IsNotExist(err) {
		os.MkdirAll(conf.Conf.DataDir+"/article", 0755)
	}

	if _, err := os.Stat(conf.Conf.DataDir + "/assets"); os.IsNotExist(err) {
		os.MkdirAll(conf.Conf.DataDir+"/assets", 0755)
	}

	if _, err := os.Stat(conf.Conf.DataDir + "/plugins"); os.IsNotExist(err) {
		os.MkdirAll(conf.Conf.DataDir+"/plugins", 0755)
	}

	if *backup {
		if err := core.Backup(); err != nil {
			log.Println(err)
		}
		os.Exit(0)
	}

	db := database.Instance()
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&article.Article{})

	gin.SetMode(conf.Conf.Mode)

	if !conf.Conf.UseCategory {
		tmp := category.Category{Name: "other"}
		if db.Where(&tmp).First(&tmp).RecordNotFound() {
			if err := db.Create(&tmp).Error; err != nil {
				log.Panicln(err)
			}
		}

		db.Model(&article.Article{}).Updates(article.Article{CategoryId: tmp.ID})
		conf.Conf.OtherCategoryId = tmp.ID
	}

	plugin.InitPlugins(loginMiddleware)
	plugin.InitDatabases()
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
	config := conf.Config{}
	fmt.Println("======================")
	fmt.Println("首次运行需要部署一些东西")
	fmt.Println("======================")
	fmt.Print("请输入博客名:")
	fmt.Scanln(&config.Name)
	fmt.Print("请输入登录密码:")
	fmt.Scanln(&config.Password)
	md5Data := md5.Sum([]byte(config.Password))
	sha1Data := sha1.Sum([]byte(md5Data[:]))
	config.Password = hex.EncodeToString(sha1Data[:])
	fmt.Println("是否开启分类功能(y/n)")
	tmp := "n"
	fmt.Scanln(&tmp)
	if tmp == "y" {
		config.UseCategory = true
	}
	fmt.Println("======================")
	fmt.Println("数据库设置")
	fmt.Println("======================")
	fmt.Println("选择数据库驱动")
	fmt.Println("1.mysql")
	fmt.Println("2.postgreSQL")
	fmt.Println("3.sqlite3")
	fmt.Println("4.mssql (SQL Server)")
	fmt.Scanln(&tmp)
	if tmp == "1" {
		config.Db.Driver = "mysql"
	} else if tmp == "2" {
		config.Db.Driver = "postgres"
	} else if tmp == "3" {
		config.Db.Driver = "sqlite3"
	} else {
		config.Db.Driver = "mssql"
	}
	fmt.Println("数据库用户名(root)")
	fmt.Scanln(&config.Db.Username)
	if config.Db.Username == "" {
		config.Db.Username = "root"
	}
	fmt.Println("数据库密码")
	fmt.Scanln(&config.Db.Password)
	fmt.Println("数据库地址(localhost)")
	fmt.Scanln(&config.Db.Address)
	if config.Db.Address == "" {
		config.Db.Address = "localhost"
	}
	fmt.Println("数据库端口(3306)")
	fmt.Scanln(&config.Db.Port)
	if config.Db.Port == "" {
		config.Db.Port = "3306"
	}
	fmt.Println("数据库名(goblog)")
	fmt.Scanln(&config.Db.Dbname)
	if config.Db.Dbname == "" {
		config.Db.Dbname = "goblog"
	}

	config.DataDir = "data"
	config.Mode = "release"
	conf.Conf = &config
	return conf.SaveConf()
}
