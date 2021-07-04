package executor

import (
	"errors"
	"meepo/task/tasks"
	"meepo/util"
	"strings"
)

func Exec(mt *tasks.MeepoTask) ([]byte, error) {
	util.SugarLogger.Debug("执行插件任务")
	b, err := mt.EncodeTask()
	if err != nil {
		util.SugarLogger.Debug("任务编码失败：%s", err.Error())
	}
	util.SugarLogger.Debugf("任务内容：%s", b)

	err = PP.LoadPlugin(mt.PluginName)
	if err != nil {
		util.SugarLogger.Error(err)
	}
	var fx interface{}
	util.SugarLogger.Debug("使用插件...")
	p := PP.Plugins[mt.PluginName].Plugin
	util.SugarLogger.Debug("查找命令...")
	fx, err = p.Lookup(strings.Title(mt.Command))
	if f, ok := fx.(func(*tasks.Args, []byte) ([]byte, error)); ok {
		util.SugarLogger.Debug("获取到命令函数")
		util.SugarLogger.Debugf("执行命令，args：\n%#v", mt.Args)
		out, err := f(mt.Args, mt.ExtraArgs)
		if err != nil {
			util.SugarLogger.Error(err)
			return nil, err
		}
		return out, nil
	}
	util.SugarLogger.Errorf("无法找到命令函数：%s", strings.Title(mt.Command))
	return nil, errors.New("can not find func name: " + mt.Command)
}
