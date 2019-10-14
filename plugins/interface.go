package plugins

import "github.com/gin-gonic/gin"

type Plugins interface {
	InitPlugins(server *PluginServer) error
	InitDatabase() error
	InitRouter(router *gin.RouterGroup) error
	GetPluginName() string
	GetPluginVersion() string
}
