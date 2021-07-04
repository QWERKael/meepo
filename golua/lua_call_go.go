package golua

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"github.com/yuin/gopher-lua"
	"meepo/common"
	"meepo/pb"
	"meepo/task/executor"
	"meepo/task/tasks"
	"meepo/util"
	"time"
	"utility-go/codec"
)

//插件运行器，在Lua中调用Go插件
func (lvm *LuaVM) PluginRunner(L *lua.LState) int {
	//pluginName := L.ToString(1) /* get argument */
	//funcName := L.ToString(2)   /* get argument */
	in := L.ToString(1) /* get argument */
	mt := &tasks.MeepoTask{}
	util.SugarLogger.Debugf("获取取到输入：%s", in)
	if err := mt.DecodeTask([]byte(in)); err != nil {
		util.SugarLogger.Errorf("解析任务失败:\n%s", err.Error())
	}
	util.SugarLogger.Debugf("识别到meepo task：%#v\nargs：%#v", mt, mt.Args)
	out, err := executor.Exec(mt)
	if err != nil {
		util.SugarLogger.Error(err)
	}
	L.Push(lua.LString(fmt.Sprintf("%s", out))) /* push result */
	return 1                                    /* number of results */
}

//插件运行器，在Lua中调用Go插件，只提供ExtraArgs参数
func (lvm *LuaVM) PluginRunnerEx(L *lua.LState) int {
	pluginName := L.ToString(1) /* get argument */
	command := L.ToString(2)    /* get argument */
	extraArgs := L.ToString(3)  /* get argument */
	mt := &tasks.MeepoTask{
		PluginName: pluginName,
		Command:    command,
		Args:       nil,
		ExtraArgs:  []byte(extraArgs),
		Result:     nil,
	}
	out, err := executor.Exec(mt)
	if err != nil {
		util.SugarLogger.Error(err)
	}
	L.Push(lua.LString(fmt.Sprintf("%s", out))) /* push result */
	return 1                                    /* number of results */
}

//将util.Config编码为json格式字符串，作为返回值传给lua
func (lvm *LuaVM) GetConfigJson(L *lua.LState) int {
	b, err := codec.EncodeJson(util.Config)
	if err != nil {
		util.SugarLogger.Errorf("将配置文件导出为Json格式时出错！")
	}
	L.Push(lua.LString(fmt.Sprintf("%s", b)))
	return 1
}

//发送任务到远端服务器
func (lvm *LuaVM) SendTask(L *lua.LState) int {
	remoteServerAddress := L.ToString(1) /* get argument */
	id := L.ToString(2)                  /* get argument */
	task := L.ToString(3)                /* get argument */
	remoteServerPID := actor.NewPID(remoteServerAddress, id)

	util.SugarLogger.Debugf("生成远端PID：\n%#v", remoteServerPID)

	req := &pb.Request{
		UUID:    uuid.NewString(),
		Sender:  common.Mee.LocalServerPID,
		MsgType: pb.MEEPO_TASK,
		ByteMsg: []byte(task),
	}

	rstChan := make(chan []byte, 1)
	tasks.ResultCache.Add(req.UUID, rstChan)

	common.Mee.RootContext.Send(remoteServerPID, req)

	L.Push(lua.LString(req.UUID))
	return 1 /* number of results */
}

//发送任务到远端服务器，同步
func (lvm *LuaVM) SendTaskRealTime(L *lua.LState) int {
	remoteServerAddress := L.ToString(1) /* get argument */
	id := L.ToString(2)                  /* get argument */
	task := L.ToString(3)                /* get argument */
	remoteServerPID := actor.NewPID(remoteServerAddress, id)

	util.SugarLogger.Debugf("生成远端PID：\n%#v", remoteServerPID)

	req := &pb.Request{
		UUID:    uuid.NewString(),
		Sender:  common.Mee.LocalServerPID,
		MsgType: pb.MEEPO_TASK,
		ByteMsg: []byte(task),
	}

	rstChan := make(chan []byte, 1)
	tasks.ResultCache.Add(req.UUID, rstChan)

	fut := common.Mee.RootContext.RequestFuture(remoteServerPID, req, 30 * time.Second)
	rst, err := fut.Result()
	if err != nil {
		util.SugarLogger.Debugf("Future 出错：\n%s", err.Error())
	}

	L.Push(lua.LString(fmt.Sprintf("%s", rst)))
	return 1 /* number of results */
}

//从缓存中根据uuid同步的返回执行结果
func (lvm *LuaVM) GetResult(L *lua.LState) int {
	UUID := L.ToString(1) /* get argument */
	rst := tasks.GetResult(UUID)

	util.SugarLogger.Debugf("获取到rst：%s", rst)
	L.Push(lua.LString(fmt.Sprintf("%s", rst)))
	return 1 /* number of results */
}

//休眠
func (lvm *LuaVM) Sleep(L *lua.LState) int {
	//sec := L.ToInt64(1) /* get argument */

	util.SugarLogger.Debugf("休眠前")
	time.Sleep(3 * time.Second)
	util.SugarLogger.Debugf("休眠后")
	//L.Push(lua.LString(fmt.Sprintf("%s", rst)))
	return 0 /* number of results */
}
