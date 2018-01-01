package api

import (
	"fmt"
	"crypto/md5"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/1377195627/goblog/token"
	"net/http"
	"github.com/1377195627/goblog"
	"log"
	"time"
	"io/ioutil"
	"github.com/russross/blackfriday"
)

type articleData struct {
	OldName string `json:"oldName" form:"oldName"`
	Name    string `json:"name" form:"name"`
	Tag     string `json:"tag" form:"tag"`
	Context string `json:"context" form:"context"`
}

func Login(context *gin.Context) {
	passwd := context.DefaultPostForm("password", "")
	passwd = fmt.Sprintf("%x", md5.Sum([]byte(passwd)))
	passwd = fmt.Sprintf("%x", md5.Sum([]byte(passwd)))

	if goblog.Conf.Password != passwd {
		context.JSON(http.StatusOK, gin.H{
			"token": "",
			"ttl":   "",
		})
	} else {

		t := token.GetManager().GetToken()

		context.JSON(http.StatusOK, gin.H{
			"token": t.Token,
			"ttl":   t.Ttl,
		})
	}

}

func Name(context *gin.Context)  {
	context.String(http.StatusOK,goblog.Conf.Name)
}

func Tags(context *gin.Context) {
	rows, err := goblog.DB.Query("SELECT DISTINCT tag FROM article")
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
}

func Tag(context *gin.Context) {
	rows, err := goblog.DB.Query("SELECT id,name,uuid,tag,create_time,edit_time FROM article")

	if err != nil {
		log.Fatal(err)
	}

	var articles []goblog.Article

	for rows.Next() {
		var article goblog.Article
		var createTime, editTIme string
		rows.Scan(&article.Id,&article.Name, &article.Uuid, &article.Tag, &createTime, &editTIme)
		t, _ := time.Parse("2006-01-02 15:04:05", createTime)
		article.CreateTime = t.Unix()
		t, _ = time.Parse("2006-01-02 15:04:05", editTIme)
		article.EditTime = t.Unix()
		articles = append(articles, article)
	}

	context.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func TagBytag(context *gin.Context) {
	tag := context.Param("tag")

	rows, err := goblog.DB.Query("SELECT id,name,uuid,tag,create_time,edit_time FROM article WHERE tag=?", tag)
	if err != nil {
		log.Fatal(err)
	}

	var articles []goblog.Article

	for rows.Next() {
		var article goblog.Article
		var createTime, editTIme string
		rows.Scan(&article.Id,&article.Name, &article.Uuid, &article.Tag, &createTime, &editTIme)
		t, _ := time.Parse("2006-01-02 15:04:05", createTime)
		article.CreateTime = t.Unix()
		t, _ = time.Parse("2006-01-02 15:04:05", editTIme)
		article.EditTime = t.Unix()
		articles = append(articles, article)
	}

	context.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func ArticleByName(context *gin.Context) {
	name := context.Param("name")

	row := goblog.DB.QueryRow("SELECT id,name,uuid,tag,create_time,edit_time FROM article WHERE name=?", name)

	var article goblog.Article
	var createTime, editTIme string
	err := row.Scan(&article.Id,&article.Name, &article.Uuid, &article.Tag, &createTime, &editTIme)
	t, _ := time.Parse("2006-01-02 15:04:05", createTime)
	article.CreateTime = t.Unix()
	t, _ = time.Parse("2006-01-02 15:04:05", editTIme)
	article.EditTime = t.Unix()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"article": article,
			"status":  err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"article": article,
		"status":  "success",
	})
}

func ArticleByUuid(context *gin.Context) {
	mode := context.DefaultQuery("mode", "complete")
	md, err := ioutil.ReadFile("article/" + context.Param("uuid") + ".md")
	if err != nil {
		log.Fatal(err)
		context.String(http.StatusNotFound, "")
	}

	if mode == "raw" {
		context.String(http.StatusOK, string(md))
	} else if mode == "complete" {
		html := blackfriday.MarkdownBasic(md)
		context.String(http.StatusOK, string(html))
	} else {
		html := blackfriday.MarkdownBasic(md)
		if len(html) > 100 {
			context.String(http.StatusOK, string(html[:99]))
		} else {
			context.String(http.StatusOK, string(html))
		}
	}
}

func ArticleNew(context *gin.Context) {
	t := context.PostForm("token")
	if !token.GetManager().IsExist(t) {
		context.JSON(http.StatusOK, gin.H{
			"status": "no authorized",
		})
		return
	}
	var data articleData
	context.Bind(&data)
	if err := goblog.AddArticle(data.Name, data.Tag, data.Context); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func ArticleDel(context *gin.Context) {
	t := context.Query("token")
	if !token.GetManager().IsExist(t) {
		context.JSON(http.StatusOK, gin.H{
			"status": "no authorized",
		})
		return
	}
	name := context.Param("name")
	if err := goblog.DelArticle(name); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func ArticleEdit(context *gin.Context) {
	t := context.PostForm("token")
	if !token.GetManager().IsExist(t) {
		context.JSON(http.StatusOK, gin.H{
			"status": "no authorized",
		})
		return
	}
	var data articleData
	context.Bind(&data)
	fmt.Println(data)
	if err := goblog.UpdateArticle(data.OldName, data.Name, data.Tag, data.Context); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
