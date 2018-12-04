package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/article"
	"gitlab.com/xiayesuifeng/goblog/database"
)

type Tag struct {
}

func (t *Tag) Get(ctx *gin.Context) {
	db := database.Instance()
	articles := make([]article.Article, 0)

	db.Order("created_at DESC").Find(&articles).Where("tag", ctx.Param("tag"))

	ctx.JSON(200, gin.H{
		"code":     0,
		"articles": articles,
	})
}

func (t *Tag) Gets(ctx *gin.Context) {
	db := database.Instance()
	rows, err := db.Table("articles").Select("DISTINCT tag").Rows()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	tags := make([]string, 0)

	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}

	ctx.JSON(200, gin.H{
		"code": 0,
		"tags": tags,
	})
}
