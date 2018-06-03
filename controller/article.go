package controller

import (
	"github.com/1377195627/goblog/article"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/1377195627/goblog/database"
	"io/ioutil"
	"github.com/russross/blackfriday"
)

type Article struct {
}

func (a *Article) Get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      100,
			"message": "id must integer",
		})
		return
	}

	db := database.Instance()
	article := article.Article{}

	db.First(&article,id)

	ctx.JSON(200,gin.H{
		"code":0,
		"article":article,
	})
}

func (a *Article) GetByCategory(ctx *gin.Context) {
	param := ctx.Param("category_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      100,
			"message": "id must integer",
		})
		return
	}

	db := database.Instance()
	articles := make([]article.Article,0)

	db.Find(&articles).Where("category_id",id)

	ctx.JSON(200,gin.H{
		"code":0,
		"articles":articles,
	})
}
func (a *Article) GetByUuid(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	mode := ctx.Param("mode")

	md, err := ioutil.ReadFile("data/article/" + uuid + ".md")
	if err != nil {
		ctx.JSON(200,gin.H{
			"code":100,
			"message":"uuid not found",
		})
		return
	}

	switch mode {
	case "description":
		html := blackfriday.MarkdownBasic(md)
		data := []rune(string(html))
		ctx.JSON(200,gin.H{
			"code":0,
			"html":string(data[:100]),
		})
	case "html":
		html := blackfriday.MarkdownBasic(md)
		ctx.JSON(200,gin.H{
			"code":0,
			"html":string(html),
		})
	case "markdown":
		ctx.JSON(200,gin.H{
			"code":0,
			"markdown":string(md),
		})
	}
}

func (a *Article) Gets(ctx *gin.Context) {
	articles := make([]article.Article,0)
	db := database.Instance()
	db.Find(&articles)

	ctx.JSON(200,gin.H{
		"code":0,
		"articles":articles,
	})
}

func (a *Article) Post(ctx *gin.Context) {
	var data struct {
		article.Article
		Context string `json:"context" binding:"required"`
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	if err := article.AddArticle(data.Title, data.Tag, data.CategoryId, data.Context); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": 0,
		})
	}
}

func (a *Article) Put(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      100,
			"message": "id must integer",
		})
	}

	var data struct {
		article.Article
		Context string `json:"context" binding:"required"`
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	if err:=article.EditArticle(uint(id), data.CategoryId, data.Title, data.Tag, data.Context); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": 0,
		})
	}
}

func (a *Article) Delete(ctx *gin.Context) {
	param := ctx.Param("id")
	id,err :=strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":      100,
			"message": "id must integer",
		})
		return
	}

	if err:= article.DeleteArticle(id);err!=nil{
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	}else{
		ctx.JSON(200, gin.H{
			"code":    0,
		})
	}
}
