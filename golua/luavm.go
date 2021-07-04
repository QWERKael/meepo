package golua

import (
	"fmt"
	libs "github.com/vadv/gopher-lua-libs"
	"meepo/common"
	//"github.com/tengattack/gluasql"
	lua "github.com/yuin/gopher-lua"
	"layeh.com/gopher-lfs"
	"meepo/util"
)

var LVM LuaVM

type LuaVM struct {
	L *lua.LState
	//Serv *server.Server
}

func init() {
	fmt.Println("初始化模块: luavm.go")

	LVM = LuaVM{
		L: lua.NewState(),
	}

	common.Mee.L = LVM.L

	//defer L.Close()
	if err := LVM.SetLuaGlobalValues(); err != nil {
		util.SugarLogger.Error("can not set package path: ", err)
	}
	util.SugarLogger.Debug("add package path: ", util.Config.LuaPackagePathAbs)
}

// SetLuaGlobalValues 设置lua的全局变量
func (lvm *LuaVM) SetLuaGlobalValues() error {
	doCode := fmt.Sprintf(`
-- 将Lua的包地址加载到Lua虚拟机的系统变量中
package.path = "%s;"..package.path
-- 设置package.config值
_G.package.config = [[/
;
?
!
-
]]
`, util.Config.LuaPackagePathAbs)

	util.SugarLogger.Debug("Lua预执行命令： ", doCode)
	if err := lvm.L.DoString(doCode); err != nil {
		return err
	}
	return nil
}

// LoadAs 将Go函数加载为Lua函数
func (lvm *LuaVM) LoadAs(fn lua.LGFunction, funcName string) {
	util.SugarLogger.Debug("Load function ", funcName)
	lvm.L.SetGlobal(funcName, lvm.L.NewFunction(fn))
}

func (lvm *LuaVM) InitLoad() {
	util.SugarLogger.Debugf("开始加载初始化 lua 模块...")
	lvm.LoadAs(lvm.PluginRunner, "runner")
	lvm.LoadAs(lvm.PluginRunnerEx, "runnerEx")
	lvm.LoadAs(lvm.GetConfigJson, "getConfig")
	lvm.LoadAs(lvm.SendTask, "send")
	lvm.LoadAs(lvm.SendTaskRealTime, "sendRT")
	lvm.LoadAs(lvm.GetResult, "getRst")
	lvm.LoadAs(lvm.Sleep, "sleep")
	lfs.Preload(lvm.L)
	libs.Preload(lvm.L)
	//gluasql.Preload(lvm.L)
	util.SugarLogger.Debugf("加载初始化 lua 模块完毕！")
}

func (lvm *LuaVM) ExecuteLuaScriptWithResult(scriptPath string) (lua.LValue, error) {
	return ExecuteLuaScriptWithResult(lvm.L, scriptPath)
}

func (lvm *LuaVM) ExecuteLuaScriptWithArgsResult(scriptPath string, luaFuncName string, args ...string) (lua.LValue, error) {
	return ExecuteLuaScriptWithArgsResult(lvm.L, scriptPath, luaFuncName, args...)
}

func (lvm *LuaVM) ExecuteLuaCommand(command string) error {
	return ExecuteLuaCommand(lvm.L, command)
}
