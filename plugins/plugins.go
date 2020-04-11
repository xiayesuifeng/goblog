package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.com/xiayesuifeng/goblog/conf"
)

type PluginServer struct {
	Config *conf.Config

	LoginMiddleware func(ctx *gin.Context)

	DatabaseInstance func() *gorm.DB
}
