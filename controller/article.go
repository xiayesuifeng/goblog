package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"gitlab.com/xiayesuifeng/goblog/article"
	"gitlab.com/xiayesuifeng/goblog/conf"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/database"
	"io/ioutil"
	"strconv"
)

type Article struct {
}

func (a *Article) Get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	db := database.Instance()
	article := article.Article{}

	db.First(&article, id)

	ctx.JSON(200, core.SuccessDataResult("article", article))
}

func (a *Article) GetByCategory(ctx *gin.Context) {
	param := ctx.Param("category_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	db := database.Instance()
	articles := make([]article.Article, 0)

	db.Where("category_id = ?", id).Order("created_at DESC").Find(&articles)

	ctx.JSON(200, core.SuccessDataResult("articles", articles))
}
func (a *Article) GetByUuid(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	mode := ctx.Param("mode")

	md, err := ioutil.ReadFile("data/article/" + uuid + ".md")
	if err != nil {
		ctx.JSON(200, core.FailResult("uuid not found"))
		return
	}

	switch mode {
	case "description":
		tmp := []rune(string(md))
		if len(tmp) > 100 {
			tmp = tmp[:100]
		}
		html := blackfriday.MarkdownBasic([]byte(string(tmp)))
		ctx.JSON(200, core.SuccessDataResult("html", string(html)))
	case "description_md":
		tmp := []rune(string(md))
		if len(tmp) > 100 {
			tmp = tmp[:100]
		}
		ctx.JSON(200, core.SuccessDataResult("markdown", string(tmp)))
	case "html":
		html := blackfriday.MarkdownBasic(md)
		ctx.JSON(200, core.SuccessDataResult("html", string(html)))
	case "markdown":
		ctx.JSON(200, core.SuccessDataResult("markdown", string(md)))
	}
}

func (a *Article) Gets(ctx *gin.Context) {
	articles := make([]article.Article, 0)
	db := database.Instance()
	db.Order("created_at DESC").Find(&articles)

	ctx.JSON(200, core.SuccessDataResult("articles", articles))
}

func (a *Article) Post(ctx *gin.Context) {
	var data struct {
		article.Article
		Context string `json:"context" binding:"required"`
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
		return
	}

	if conf.Conf.UseCategory {
		if data.CategoryId == 0 {
			ctx.JSON(200, core.FailResult("category_id must exist"))
			return
		}
	} else {
		data.CategoryId = conf.Conf.OtherCategoryId
	}

	if err := article.AddArticle(data.Title, data.Tag, data.CategoryId, data.Context); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}

func (a *Article) Put(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
	}

	var data struct {
		article.Article
		Context string `json:"context" binding:"required"`
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
		return
	}

	if conf.Conf.UseCategory {
		if data.CategoryId == 0 {
			ctx.JSON(200, core.FailResult("category_id must exist"))
			return
		}
	} else {
		data.CategoryId = conf.Conf.OtherCategoryId
	}

	if err := article.EditArticle(uint(id), data.CategoryId, data.Title, data.Tag, data.Context); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}

func (a *Article) Delete(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	if err := article.DeleteArticle(id); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}
