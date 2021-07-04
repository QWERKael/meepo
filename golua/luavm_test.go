package golua

import (
	"fmt"
	"meepo/task/executor"
	"meepo/util"
	"testing"
)

func init() {
	//设置工作目录
	util.Config.WorkDir = ".."
	err := util.Config.Default()
	if err != nil {
		util.SugarLogger.Error(err.Error())
	}
	util.SugarLogger.Debugf("配置文件：\n :%#v", util.Config)
	//重新加载lua包地址
	if err := LVM.SetLuaGlobalValues(); err != nil {
		util.SugarLogger.Error("can not set package path: ", err)
	}
	//重新初始化PluginPool
	executor.PluginPoolInit()
}

func TestPluginRunner(t *testing.T) {
	defer LVM.L.Close()


	//if err := LVM.L.DoFile("../lua/hello.lua"); err != nil {
	//	panic(err)
	//}

	LVM.InitLoad()
	//if err := LVM.L.DoFile("../lua/runner.lua"); err != nil {
	//	panic(err)
	//}

	fmt.Println("***************************************************************")

	if err := LVM.L.DoFile("../lua/sender.lua"); err != nil {
		panic(err)
	}
	// 执行字符串语句
	//if err := LVM.L.DoString(`print("HELLO WORLD!")`); err != nil {
	//	panic(err)
	//}
	//
	//if rst, err := ExecuteLuaScriptWithResult(LVM.L, "../lua/funcs.lua"); err != nil {
	//	panic(err)
	//} else {
	//	fmt.Printf("a + b = %#v", rst)
	//	fmt.Println()
	//	rst.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {fmt.Println(k, ":", v)})
	//}
}