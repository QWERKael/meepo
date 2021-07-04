package server

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	myactor "meepo/actor"
	"meepo/common"
	"meepo/golua"
	"meepo/pb"
	"meepo/task/executor"
	"meepo/task/tasks"
	"meepo/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Serv = &Server{}

type Server struct {
	RootContext    *actor.RootContext
	LocalServerPID *actor.PID
}

func (s *Server) StartServer() error {
	rootContext := myactor.NewRootContext(util.Config.ListenHost, util.Config.ListenPort)
	s.RootContext = rootContext
	common.Mee.RootContext = s.RootContext
	fileMap := make(map[string]*common.FileBuilder)
	props := actor.PropsFromFunc(
		func(context actor.Context) {
			util.SugarLogger.Debugf("!!!!! Sender: %#v", context.Sender())
			switch msg := context.Message().(type) {
			case *pb.Request:
				go func(msg *pb.Request) {
					mt := &tasks.MeepoTask{}
					if err := mt.DecodeTask(msg.ByteMsg); err != nil {
						util.SugarLogger.Errorf("解析任务失败:\n%s", err.Error())
					}

					util.SugarLogger.Debugf("捕获到【%s】的请求:\n%#v", msg.UUID, mt)
					var (
						rst   []byte
						state pb.Response_StateCode
					)

					switch mt.PluginName {
					case "lua":
						golua.LVM.InitLoad()
						err := golua.LVM.ExecuteLuaCommand(mt.Command)
						if err != nil {
							rst = []byte("执行失败：" + err.Error())
						} else {
							rst = []byte("执行成功")
						}
					case "lua-script":
						golua.LVM.InitLoad()
						scriptPath := filepath.Join(util.Config.LuaScriptsPathAbs, mt.Command)
						//lsRst, err := golua.LVM.ExecuteLuaScriptWithResult(scriptPath)
						var (
							luaFuncName string
							subCommands []string
						)

						if mt.Args != nil {
							if v, ok := mt.Args.KVs["func"]; ok {
								luaFuncName = v[0]
							} else {
								luaFuncName = "main"
							}
							subCommands = mt.Args.SubCommands
						} else {
							luaFuncName = "main"
						}
						util.SugarLogger.Debugf("执行lua脚本【%s】的【%s】函数，参数为：%#v", scriptPath, luaFuncName, subCommands)
						lsRst, err := golua.LVM.ExecuteLuaScriptWithArgsResult(scriptPath, luaFuncName, subCommands...)
						if err != nil {
							rst = []byte("执行失败：" + err.Error())
						} else {
							rst = []byte(fmt.Sprintf("%#v", lsRst))
						}
					default:
						var err error
						rst, err = executor.Exec(mt)
						if err != nil {
							rst = []byte("执行失败：" + err.Error())
						}
					}
					util.SugarLogger.Debugf("返回给【%s】的结果:\n%s", msg.Sender.Address, rst)
					rootContext.Send(msg.Sender, &pb.Response{
						UUID:    msg.UUID,
						State:   state,
						ByteMsg: rst,
					})
				}(msg)

			case *pb.Response:
				go func(msg *pb.Response) {
					if val, ok := tasks.ResultCache.Get(msg.UUID); ok {
						if rstChan, ok := val.(chan []byte); ok {
							rstChan <- msg.ByteMsg
							close(rstChan)
						} else {
							util.SugarLogger.Errorf("获取到UUID【%s】的结果不是【chan []byte】类型", msg.UUID)
						}
					} else {
						util.SugarLogger.Errorf("获取不到UUID【%s】的结果", msg.UUID)
					}
					util.SugarLogger.Infof("获取到【%s】的结果:\n%s", msg.UUID, msg.ByteMsg)
				}(msg)

			case *pb.Transfer:
				//util.SugarLogger.Debugf("接收到文件传输请求：%#v", msg)
				//go func(msg *pb.Transfer) {
				//	switch msg.State {
				//	case pb.Req:
				//		util.SugarLogger.Debugf("接收到文件传输Req请求")
				//		ctx := string(msg.Context)
				//		util.SugarLogger.Debugf("接收到文件传输请求【%s】", ctx)
				//		c := strings.Split(ctx, "|")[:2]
				//		fileBuilder := &common.FileBuilder{
				//			FileUUID: msg.FileUUID,
				//			Sender:   msg.Sender,
				//			Name:     c[0],
				//			MD5:      c[1],
				//			File:     nil,
				//			Channel:  make(chan []byte),
				//			Done:     false,
				//		}
				//		var err error
				//		filePath := path.Join(util.Config.WorkDirAbs, "files", fileBuilder.Name)
				//		if fileBuilder.File, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
				//			util.SugarLogger.Errorf("写文件【%s】失败：%s", filePath, err.Error())
				//		}
				//		fileMap[msg.FileUUID] = fileBuilder
				//	case pb.Done:
				//		util.SugarLogger.Debugf("接收到文件传输Done请求")
				//		fallthrough
				//	case pb.Data:
				//		util.SugarLogger.Debugf("接收到文件传输Data请求")
				//		fileBuilder,ok := fileMap[msg.FileUUID]
				//		if !ok{
				//			return
				//		}
				//		//err := ioutil.WriteFile(path.Join(util.Config.WorkDirAbs, "files", fileMap[msg.FileUUID].Name), msg.Context, 0644)
				//		//if err != nil {
				//		//	util.SugarLogger.Errorf("文件【%s】传输失败", fileMap[msg.FileUUID].Name)
				//		//}
				//		var err error
				//		if !fileBuilder.Done {
				//			util.SugarLogger.Debugf("接受文件【%s】的数据", fileBuilder.Name)
				//			b := msg.Context
				//			util.SugarLogger.Debugf("接受到文件【%s】的数据", fileBuilder.Name)
				//			n, err := fileBuilder.File.Write(b[:])
				//			if err != nil {
				//				util.SugarLogger.Errorf("文件【%s】写入出错：%s", fileBuilder.Name, err.Error())
				//				return
				//			}
				//			util.SugarLogger.Debugf("文件【%s】写入数据 %d 字节", fileBuilder.Name, n)
				//		}
				//		if msg.State == pb.Done {
				//			err = fileBuilder.File.Close()
				//			if err != nil {
				//				util.SugarLogger.Errorf("无法关闭文件【%s】：%s！", fileBuilder.Name, err.Error())
				//			}
				//		}
				//	}
				//
				//}(msg)

				// ===========================================================================================================================
				util.SugarLogger.Debugf("接收到文件传输请求：%#v", msg)
				go func() {
					switch msg.State {
					case pb.Req:
						ctx := string(msg.Context)
						util.SugarLogger.Debugf("接收到文件传输请求【%s】", ctx)
						c := strings.Split(ctx, "|")[:2]
						fileBuilder := &common.FileBuilder{
							FileUUID: msg.FileUUID,
							Sender:   msg.Sender,
							Name:     c[0],
							MD5:      c[1],
							File:     nil,
							Channel:  make(chan []byte),
							Done:     false,
						}
						filePath := path.Join(util.Config.WorkDirAbs, "files", fileBuilder.Name)
						var err error
						fileMap[msg.FileUUID] = fileBuilder

						if fileBuilder.File, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
							util.SugarLogger.Errorf("写文件【%s】失败：%s", filePath, err.Error())
							context.Request(msg.Sender, &pb.Transfer{
								FileUUID: msg.FileUUID,
								Sender:   s.LocalServerPID,
								State:    pb.Refuse,
								Size_:    0,
								Context:  nil,
							})
						} else {
							go func(fileBuilder *common.FileBuilder) {
								for !fileBuilder.Done {
									util.SugarLogger.Debugf("接受文件【%s】的数据", fileBuilder.Name)
									b := <-fileBuilder.Channel
									util.SugarLogger.Debugf("接受到文件【%s】的数据", fileBuilder.Name)
									n, err := fileBuilder.File.Write(b[:])
									if err != nil {
										util.SugarLogger.Errorf("文件【%s】写入出错：%s", fileBuilder.Name, err.Error())
										context.Request(fileBuilder.Sender, &pb.Transfer{
											FileUUID: msg.FileUUID,
											Sender:   s.LocalServerPID,
											State:    pb.Fail,
											Size_:    0,
											Context:  nil,
										})
										return
									}
									util.SugarLogger.Debugf("文件【%s】写入数据 %d 字节", fileBuilder.Name, n)
								}
								err = fileBuilder.File.Close()
								if err != nil {
									util.SugarLogger.Errorf("无法关闭文件【%s】：%s！", filePath, err.Error())
								}
							}(fileBuilder)

							context.Respond(&pb.Transfer{
								FileUUID: msg.FileUUID,
								Sender:   s.LocalServerPID,
								State:    pb.OK,
								Size_:    0,
								Context:  nil,
							})
							util.SugarLogger.Debugf("Sender: %#v", context.Sender())

							util.SugarLogger.Debugf("文件【%s】已创建，准备接收数据", fileBuilder.Name)
						}
						return
					case pb.Data:
						fileMap[msg.FileUUID].Channel <- msg.Context
					case pb.Done:
						fileMap[msg.FileUUID].Channel <- msg.Context
						fileMap[msg.FileUUID].Done = true
					case pb.Close:
						fileMap[msg.FileUUID].Done = true
					}
				}()

				// ===========================================================================================================================

			}
		})
	var err error
	s.LocalServerPID, err = rootContext.SpawnNamed(props, util.Config.ActorName)
	common.Mee.LocalServerPID = s.LocalServerPID
	if err != nil {
		util.SugarLogger.Errorf("Spawn server error: %s", err.Error())
		return err
	}
	return nil
}
