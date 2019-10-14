package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.com/xiayesuifeng/goblog/core"
	"plugin"
)

type PluginServer struct {
	Config *core.Config

	LoginMiddleware func(ctx *gin.Context)

	DatabaseInstance func() *gorm.DB
}

type PluginClient struct {
	*plugin.Plugin

	Plugins

	pluginName string
}
