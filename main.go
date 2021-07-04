package main

import (
	"context"
	"meepo/golua"
	"meepo/node/client"
	"meepo/node/server"
	"meepo/util"
)

func main() {
	defer golua.LVM.L.Close()
	if err := util.ParseConfig(); err != nil {
		util.SugarLogger.Fatalf("init config error: %s\n", err.Error())
	}
	util.SugarLogger.Debugf("%#v", util.Config)

	//L := golua.L
	//defer L.Close()
	//if err := L.DoFile("lua/hello.lua"); err != nil {
	//	panic(err)
	//}
	//
	//util.SugarLogger.Debug("PluginName runner is running...")
	//L.SetGlobal("runner", L.NewFunction(golua.PluginRunner))
	//if err := L.DoFile("lua/runner.lua"); err != nil {
	//	panic(err)
	//}

	if err := server.Serv.StartServer();err != nil {
		util.SugarLogger.Fatalf("启动服务时发生错误: %s\n", err.Error())
	}
	//var remoteServerPID *actor.PID
	ctx, cancel := context.WithCancel(context.Background())
	//ctx = metadata.NewOutgoingContext(ctx,
	//	metadata.Pairs(
	//		"name", *name,
	//		"auth", auth,
	//	),
	//)

	//cons := client.CreateClient(remoteServerPID, localServerPID, rootContext)
	//cons.Run()

	//golua.LVM.Serv = server.Serv
	//common.Mee.Serv = server.Serv

	go func() {
		cli := &client.Client{
			Prompt:          nil,
			RemoteServerPID: server.Serv.LocalServerPID,
			LocalServerPID:  server.Serv.LocalServerPID,
			Preload:         "",
			Ctx:             ctx,
			Cancel:          cancel,
			RootContext:     server.Serv.RootContext,
		}

		client.LoadPrompt(util.Config.PromptConfigPathAbs)

		cli.Prepare()
		cli.Prompt.Run()
	} ()

	<-ctx.Done()
	util.SugarLogger.Debugf("Bye")
}
