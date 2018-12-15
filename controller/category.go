package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/category"
	"gitlab.com/xiayesuifeng/goblog/core"
	"strconv"
)

type Category struct {
}

func (c *Category) Post(ctx *gin.Context) {
	data := category.Category{}
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, core.FailResult("name is null"))
		return
	}

	if err := category.AddCategory(data.Name); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}

func (c *Category) Gets(ctx *gin.Context) {
	ctx.JSON(200, core.SuccessDataResult("categorys", category.GetCategorys()))
}

func (c *Category) Get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	category, err := category.GetCategory(uint(id))
	if err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessDataResult("category", category))
	}
}

func (c *Category) Put(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	data, err := category.GetCategory(uint(id))
	if err != nil {
		ctx.JSON(200, core.FailResult("category not found"))
		return
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(200, core.FailResult("name is null"))
		return
	}

	if err := data.SetName(data.Name); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}

func (c *Category) Delete(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(200, core.FailResult("id must integer"))
		return
	}

	if err := category.DeleteCategory(id); err != nil {
		ctx.JSON(200, core.FailResult(err.Error()))
	} else {
		ctx.JSON(200, core.SuccessResult())
	}
}
