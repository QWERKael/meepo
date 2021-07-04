package executor

import (
	"meepo/util"
	"path/filepath"
	"plugin"
	"strings"
)

var PP *PluginPool

func init() {
	PluginPoolInit()
}

func PluginPoolInit() {
	PP = &PluginPool{
		PluginDir: util.Config.PluginDirAbs,
		Plugins: make(map[string]PlugInfo),
	}
}

type PluginState uint8

const (
	Unavailable PluginState = 0
	Available   PluginState = 1
	Replaced    PluginState = 2
	Removed     PluginState = 3
)

type PlugInfo struct {
	Name   string
	Path   string
	State  PluginState
	Plugin *plugin.Plugin
}

type PluginPool struct {
	PluginDir string
	Plugins   map[string]PlugInfo
}

func (pp *PluginPool) LoadPlugin(pluginName string) error {
	if pi, ok := pp.Plugins[pluginName]; ok {
		if pi.State == Available {
			util.SugarLogger.Debugf("插件已存在")
			return nil
		}
	}
	pluginPath := filepath.Join(pp.PluginDir, strings.ToLower(pluginName)+".so")
	util.SugarLogger.Debug("加载插件", pluginPath, "中...")
	p, err := plugin.Open(pluginPath)
	if err != nil {
		util.SugarLogger.Error(err)
		return err
	}
	pluginInfo := PlugInfo{Name: pluginName, Path: pluginPath, State: Available, Plugin: p}
	pp.Plugins[pluginName] = pluginInfo
	util.SugarLogger.Debug("插件已被插入管理器")
	return nil
}