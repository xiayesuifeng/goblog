package controller

import (
	"github.com/1377195627/goblog/article"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Article struct {
}

func (a *Article) Get(ctx *gin.Context) {

}

func (a *Article) GetByCategory(ctx *gin.Context) {

}
func (a *Article) GetByUuid(ctx *gin.Context) {

}

func (a *Article) Gets(ctx *gin.Context) {

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
