package plugins

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/database"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"strings"
)

var pluginServer *PluginServer
var pluginList = make(map[string]PluginClient)

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

func InitPlugins(loginMiddleware func(ctx *gin.Context)) {
	pluginServer = &PluginServer{}

	pluginServer.Config = core.Conf
	pluginServer.LoginMiddleware = loginMiddleware
	pluginServer.DatabaseInstance = database.Instance

	file, err := ioutil.ReadDir("plugins")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("[Plugins] plugins directory not exist, skip plugin loading")
		}
		return
	}

	for _, info := range file {
		if !info.IsDir() {
			name := info.Name()
			if strings.HasSuffix(name, ".so") {
				tmp := name[:strings.Index(name, ".so")]

				log.Println("[Plugins] loading plugin: " + tmp)
				if err := registerPlugin(tmp); err != nil {
					log.Printf("[Plugins] [%s] failed to load! error: %s", tmp, err.Error())
				} else {
					log.Printf("[Plugins] [%s] loading plugin success", tmp)
				}
			}
		}
	}
}

func InitDatabases() {
	for name, pluginClient := range pluginList {
		log.Printf("[Plugins] [%s] initialize plugin database", name)
		if err := pluginClient.InitDatabase(); err != nil {
			log.Printf("[Plugins] [%s] failed to initialize database! error: %s", name, err.Error())
		} else {
			log.Printf("[Plugins] [%s] initialize database success", name)
		}
	}
}

func InitRouters(router *gin.RouterGroup) {
	for name, pluginClient := range pluginList {
		log.Printf("[Plugins] [%s] initialize plugin router", name)
		if err := pluginClient.InitRouter(router); err != nil {
			log.Printf("[Plugins] [%s] failed to initialize router! error: %s", name, err.Error())
		} else {
			log.Printf("[Plugins] [%s] initialize plugin router success", name)
		}
	}
}

func registerPlugin(name string) error {
	p, err := plugin.Open("plugins/" + name + ".so")
	if err != nil {
		return err
	}

	symbol, err := p.Lookup(strings.Title(name + "Plugin"))
	if err != nil {
		return err
	}

	pc := PluginClient{Plugin: p, pluginName: name, Plugins: symbol.(Plugins)}

	if err := pc.InitPlugins(pluginServer); err != nil {
		return err
	}

	pluginList[name] = pc

	return nil
}

func GetPluginNameList() []string {
	list := make([]string, 0)

	for name, _ := range pluginList {
		list = append(list, name)
	}

	return list
}

func GetPlugin(name string) (PluginClient, error) {
	if p, ok := pluginList[name]; ok {
		return p, nil
	} else {
		return p, errors.New("plugin not found")
	}
}
