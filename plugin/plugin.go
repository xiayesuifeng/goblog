package plugin

import "gitlab.com/xiayesuifeng/goblog/plugins"

type Plugin struct {
	PackageName   string `json:"packageName"`
	PluginName    string `json:"pluginName"`
	PluginVersion string `json:"pluginVersion"`
}

func GetPlugins() []Plugin {
	list := make([]Plugin, 0)
	names := plugins.GetPluginNameList()
	for _, name := range names {
		if p, err := plugins.GetPlugin(name); err == nil {
			list = append(list, Plugin{
				PackageName:   name,
				PluginName:    p.GetPluginName(),
				PluginVersion: p.GetPluginVersion(),
			})
		}
	}

	return list
}
