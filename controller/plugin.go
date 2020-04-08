package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/plugin"
)

type Plugin struct {
}

func (Plugin) Gets(ctx *gin.Context) {
	ctx.JSON(200, core.SuccessDataResult("plugins", plugin.GetPlugins()))
}
