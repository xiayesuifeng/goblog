package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/1377195627/goblog/database"
)

type Tag struct {
}

func (t *Tag) Get(ctx *gin.Context)  {

}

func (t *Tag) Gets(ctx *gin.Context)  {
	db := database.Instance()
	rows,err :=db.Table("articles").Select("tag").Rows()
	if err != nil {
		ctx.JSON(200,gin.H{
			"code":100,
			"message":err.Error(),
		})
		return
	}

	tags := make([]string,0)

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
