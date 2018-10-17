package controller

import (
	"gitlab.com/xiayesuifeng/goblog/category"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Category struct {
}

func (c *Category) Post(ctx *gin.Context) {
	data := category.Category{}
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "name is null",
		})
		return
	}

	if err := category.AddCategory(data.Name); err != nil {
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

func (c *Category) Gets(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":      0,
		"categorys": category.GetCategorys(),
	})
}

func (c *Category) Get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "id must integer",
		})
		return
	}

	category, err := category.GetCategory(uint(id))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":     0,
			"category": category,
		})
	}
}

func (c *Category) Put(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "id must integer",
		})
		return
	}

	data, err := category.GetCategory(uint(id))
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "category not found",
		})
		return
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "name is null",
		})
		return
	}

	if err := data.SetName(data.Name); err != nil {
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

func (c *Category) Delete(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "id must integer",
		})
		return
	}

	if err := category.DeleteCategory(id); err != nil {
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
