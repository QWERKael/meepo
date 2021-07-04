package golua

import (
	"github.com/yuin/gopher-lua"
	"meepo/util"
)

func ExecuteLuaScriptWithResult(L *lua.LState, scriptPath string) (lua.LValue, error) {
	if err := L.DoFile(scriptPath); err != nil {
		return nil, err
	}
	luaFunc := L.GetGlobal("result")
	if luaFunc.Type() != lua.LTFunction {
		util.SugarLogger.Debug("未获取到result函数！")
		return lua.LString("没有返回值！"), nil
	}
	if err := L.CallByParam(lua.P{
		Fn:      luaFunc,
		NRet:    1,
		Protect: true,
	}); err != nil {
		return nil, err
	}
	rst := L.Get(-1) // returned value
	L.Pop(1)         // remove received value
	return rst, nil
}

func ExecuteLuaScriptWithArgsResult(L *lua.LState, scriptPath string, luaFuncName string, args ...string) (lua.LValue, error) {
	//执行脚本
	if err := L.DoFile(scriptPath); err != nil {
		return nil, err
	}

	//执行main命令
	var LArgs []lua.LValue
	for _, arg := range args {
		LArgs = append(LArgs, lua.LString(arg))
	}

	luaFunc := L.GetGlobal(luaFuncName)
	if luaFunc.Type() != lua.LTFunction {
		util.SugarLogger.Debug("未获取到result函数！")
		return lua.LString("没有返回值！"), nil
	}
	util.SugarLogger.Debugf("执行lua脚本【%s】的【%s】函数，参数为：%#v", scriptPath, luaFuncName, LArgs)
	if err := L.CallByParam(lua.P{
		Fn:      luaFunc,
		NRet:    1,
		Protect: true,
	}, LArgs...); err != nil {
		return nil, err
	}
	rst := L.Get(-1) // returned value
	L.Pop(1)         // remove received value
	return rst, nil
}

func ExecuteLuaCommand(L *lua.LState, command string) error {
	if err := L.DoString(command); err != nil {
		return err
	}
	return nil
}
