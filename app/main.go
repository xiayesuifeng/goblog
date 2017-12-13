package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/1377195627/goblog"
	"github.com/1377195627/goblog/token"
	"log"
	"os"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"time"
	"crypto/md5"
)

var config goblog.Config

func main() {
	var err error
	config, err = goblog.ParseConf("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = goblog.InitSql(config)
	if err != nil {
		log.Fatal(err)
	}

	route := gin.Default()
	route.LoadHTMLGlob("view/*")

	api := route.Group("api")

	api.POST("/login", func(context *gin.Context) {
		passwd := context.DefaultPostForm("password", "")
		passwd = fmt.Sprintf("%x", md5.Sum([]byte(passwd)))
		passwd = fmt.Sprintf("%x", md5.Sum([]byte(passwd)))

		if config.Password != passwd {
			context.JSON(http.StatusOK, gin.H{
				"token": "",
				"ttl":   "",
			})
		} else {

			t:=token.GetManager().GetToken()

			context.JSON(http.StatusOK, gin.H{
				"token": t.Token,
				"ttl":   t.Ttl,
			})
		}

	})

	api.GET("/tags", func(context *gin.Context) {
		rows, err := goblog.DB.Query("SELECT tag FROM article")
		if err != nil {
			log.Fatal(err)
		}

		var tags []string

		for rows.Next() {
			var tag string
			rows.Scan(&tag)
			tags = append(tags, tag)
		}

		context.JSON(http.StatusOK, gin.H{
			"tags": tags,
		})
	})

	api.GET("/tag", func(context *gin.Context) {
		rows, err := goblog.DB.Query("SELECT * FROM article")

		if err != nil {
			log.Fatal(err)
		}

		var articles []goblog.Article

		for rows.Next() {
			var article goblog.Article
			var createTime, editTIme string;
			rows.Scan(&article.Name, &article.Uuid, &article.Tag, &createTime, &editTIme)
			t, _ := time.Parse("2006-01-02 15:04:05", createTime)
			article.CreateTime = t.Unix()
			t, _ = time.Parse("2006-01-02 15:04:05", editTIme)
			article.EditTime = t.Unix()
			articles = append(articles, article)
		}

		context.JSON(http.StatusOK, gin.H{
			"articles": articles,
		})
	})

	api.GET("/tag/:tag", func(context *gin.Context) {
		tag := context.PostForm("tag")

		rows, err := goblog.DB.Query("SELECT * FROM article WHERE tag=?", tag)
		if err != nil {
			log.Fatal(err)
		}

		var articles []goblog.Article

		for rows.Next() {
			var article goblog.Article
			rows.Scan(&article.Name, &article.Uuid, &article.Tag, &article.CreateTime, &article.EditTime)
			articles = append(articles, article)
		}

		context.JSON(http.StatusOK, gin.H{
			"articles": articles,
		})
	})

	api.GET("/article/name/:name", func(context *gin.Context) {

	})

	api.GET("/article/uuid/:uuid", func(context *gin.Context) {
		mode := context.DefaultQuery("mode", "complete")
		md, err := ioutil.ReadFile("article/" + context.Param("uuid") + ".md")
		if err != nil {
			log.Fatal(err)
			context.String(http.StatusNotFound, "")
		}
		html := blackfriday.MarkdownBasic(md)
		if mode == "complete" {
			context.String(http.StatusOK, string(html))
		} else {
			if len(html) > 100 {
				context.String(http.StatusOK, string(html[:99]))
			} else {
				context.String(http.StatusOK, string(html))
			}
		}
	})

	api.POST("/article/new", func(context *gin.Context) {
		//token:=context.PostForm("token")
		//context.c
	})

	api.POST("/article/del", func(context *gin.Context) {

	})

	api.PUT("/article/edit", func(context *gin.Context) {

	})

	route.GET("/install", func(context *gin.Context) {
		context.HTML(http.StatusOK, "install.html", gin.H{})
	})

	route.POST("/install", func(context *gin.Context) {
		err = context.Bind(&config)
		if err != nil {
			log.Fatal(err)
		}

		if err = goblog.Install(config); err != nil {
			log.Fatal(err)
			context.String(http.StatusOK, "安装失败")
		}

		context.String(http.StatusOK, "安装完成")

	})

	route.GET("/", func(context *gin.Context) {
		if _, err = os.Stat("goblog.lock"); err != nil {
			if os.IsNotExist(err) {
				context.Redirect(http.StatusMovedPermanently, "/install")
			}
		}
	})

	route.Run()
}
