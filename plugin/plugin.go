package plugin

import (
	"gitlab.com/xiayesuifeng/goblog/plugins"
	"io/ioutil"
	"strings"
)

const (
	StatusLoaded = iota + 1
	StatusNotCompatible
)

type Plugin struct {
	PackageName   string `json:"packageName"`
	PluginName    string `json:"pluginName"`
	PluginVersion string `json:"pluginVersion"`
	Status        uint   `json:"status"`
}

func GetPlugins() []Plugin {
	list := make([]Plugin, 0)

	file, err := ioutil.ReadDir("plugins")
	if err != nil {
		return list
	}

	for _, info := range file {
		if !info.IsDir() {
			name := info.Name()
			if strings.HasSuffix(name, ".so") {
				name := name[:strings.Index(name, ".so")]
				if p, err := plugins.GetPlugin(name); err == nil {
					list = append(list, Plugin{
						PackageName:   name,
						PluginName:    p.GetPluginName(),
						PluginVersion: p.GetPluginVersion(),
						Status:        StatusLoaded,
					})
				} else {
					list = append(list, Plugin{
						PackageName: name,
						Status:      StatusNotCompatible,
					})
				}
			}
		}
	}

	return list
}
