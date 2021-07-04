package client

import (
	"context"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/c-bata/go-prompt"
)

type Client struct {
	Prompt          *prompt.Prompt
	RemoteServerPID *actor.PID
	LocalServerPID  *actor.PID
	Preload         string // 预加载的命令前缀
	Ctx             context.Context
	Cancel          context.CancelFunc
	RootContext     *actor.RootContext
}


//func CreateClient(remoteServerPID *actor.PID, localServerPID *actor.PID, rootContext *actor.RootContext) *console.Console {
//	cons := console.NewConsole(func(text string) {
//		if remoteServerPID == nil {
//			util.SugarLogger.Errorf("请先建立连接！")
//			return
//		}
//
//		//rootContext.Send(remoteServerPID, &pb.Response{
//		//	UUID:    "",
//		//	Sender:  localServerPID,
//		//	MsgType: pb.TASK,
//		//	ByteMsg: []byte(text),
//		//})
//	})
//
//	cons.Command("/connect", func(addrStr string) {
//		util.SugarLogger.Debugf("连接状态[before]：\n%#v", remoteServerPID)
//		util.SugarLogger.Debugf("正在建立连接 【%s】", addrStr)
//		addr := strings.Split(addrStr, "-")
//		remoteServerPID = actor.NewPID(addr[0], addr[1])
//		util.SugarLogger.Debugf("连接状态[after]：\n%#v", remoteServerPID)
//	})
//
//	cons.Command("/run", func(runCmd string) {
//		if remoteServerPID == nil {
//			util.SugarLogger.Errorf("请先建立连接！")
//			return
//		}
//
//		cmd := strings.Split(runCmd, ".")
//
//		rootContext.Send(remoteServerPID, &pb.Request{
//			UUID:   uuid.NewString(),
//			Sender: localServerPID,
//			Task: &pb.Task{
//				PluginName: cmd[0],
//				Command:   cmd[1],
//				Args:   nil,
//			},
//		})
//	})
//
//	cons.Command("/lua", func(luaCommand string) {
//		if remoteServerPID == nil {
//			util.SugarLogger.Errorf("请先建立连接！")
//			return
//		}
//
//		rootContext.Send(remoteServerPID, &pb.Request{
//			UUID:   uuid.NewString(),
//			Sender: localServerPID,
//			Task: &pb.Task{
//				PluginName: "lua",
//				Command:   luaCommand,
//				Args:   nil,
//			},
//		})
//	})
//
//	cons.Command("/script", func(scriptName string) {
//		if remoteServerPID == nil {
//			util.SugarLogger.Errorf("请先建立连接！")
//			return
//		}
//
//		rootContext.Send(remoteServerPID, &pb.Request{
//			UUID:   uuid.NewString(),
//			Sender: localServerPID,
//			Task: &pb.Task{
//				PluginName: "lua-script",
//				Command:   scriptName,
//				Args:   nil,
//			},
//		})
//	})
//
//	return cons
//}
