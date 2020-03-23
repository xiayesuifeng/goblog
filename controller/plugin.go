package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/plugins"
)

type Plugin struct {
	Name          string `json:"name"`
	PluginName    string `json:"pluginName"`
	PluginVersion string `json:"pluginVersion"`
}

func (Plugin) Gets(ctx *gin.Context) {
	list := make([]Plugin, 0)
	names := plugins.GetPluginNameList()
	for _, name := range names {
		if p, err := plugins.GetPlugin(name); err == nil {
			list = append(list, Plugin{
				Name:          name,
				PluginName:    p.GetPluginName(),
				PluginVersion: p.GetPluginVersion(),
			})
		}
	}
	ctx.JSON(200, core.SuccessDataResult("plugins", list))
}
