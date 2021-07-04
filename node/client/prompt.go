package client

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"meepo/task/tasks"
	"meepo/util"

	"meepo/pb"
	"runtime"

	"strings"
)

//var log = logger.Logger

//var suggestion = []prompt.Suggest{
//	{Text: "show ", Description: "展示命令"},
//	{Text: "net", Description: "查看网卡信息"},
//	{Text: "load", Description: "查看负载"},
//	{Text: "processlist", Description: "查看进程"},
//	{Text: "upload ", Description: "上传文件"},
//	{Text: "change-to ", Description: "改变连接到新的服务端"},
//	{Text: "restart", Description: "重启服务端"},
//	{Text: "exit", Description: "退出"},
//}

//type Prompt struct {
//	Addr      string
//	Preload   string
//	CC        pb.CommanderClient
//	Ctx       context.Context
//	Cancel    context.CancelFunc
//	Runner    *prompt.Prompt
//	PluginDir string
//}

//func (pmt *Prompt) ConnectToAddr(reAuth bool) {
//	if reAuth {
//		auth, err := utils.Simple()
//		utils.CheckErrorPanic(err)
//		if md, ok := metadata.FromOutgoingContext(pmt.Ctx); ok {
//			md.Set("auth", auth)
//			pmt.Ctx = metadata.NewOutgoingContext(pmt.Ctx, md)
//		}
//	}
//	conn, err := grpc.Dial(pmt.Addr, grpc.WithInsecure())
//	utils.CheckErrorPanic(errors.WithMessage(err, "连接到指定地址失败: "+pmt.Addr))
//	client := pb.NewCommanderClient(conn)
//	log.Infoln("连接到", pmt.Addr, "...")
//	pmt.CC = client
//	err = es.ESs["change-to"].Add(pmt.Addr, "IP地址", Cmd)
//	if err != nil {
//		log.Debugf("IP地址添加到快捷命令失败: %s", err.Error())
//	}
//	//suggestion = append(suggestion, prompt.Suggest{Text: pmt.Addr, Description: "IP地址"})
//}

func (cli *Client) executor(line string) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err == nil {
			return
		}
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()

	if strings.TrimSpace(line) == "" {
		return
	}

	if !strings.HasPrefix(line, "$set") && !strings.HasPrefix(line, "$unset") {
		line = cli.Preload + " " + line
	}
	mt, isRT, err := Parse(line)
	if err != nil {
		util.SugarLogger.Errorf("命令行解析错误：%s", err.Error())
		return
	}
	if mt == nil {
		return
	}

	var req *pb.Request

	switch strings.ToLower(mt.PluginName) {
	case "$set":
		switch strings.ToLower(mt.Command) {
		// 设置预加载命令, 如果需要重复使用相同的命令前缀, 可以设置此值, 避免重复输入
		case "preload":
			cli.Preload = mt.Args.SubCommands[0]
		}
		return
	case "$unset":
		switch strings.ToLower(mt.Command) {
		case "preload":
			cli.Preload = ""
		}
		return
	case "get":
		b := tasks.GetResult(mt.Command)
		util.SugarLogger.Debugf("获取到UUID【%s】的返回值：\n%s", mt.Command, b)
		return
	case "":
		return
	case "exit":
		cli.Cancel()
		return
	//case "upload":
	//	mt.Type = pb.CommonCmdRequest_FILE_TRANSFER
	//case "download":
	//	mt.Type = pb.CommonCmdRequest_FILE_TRANSFER
	case "connect":
		switch mt.Command {
		case "info":
			util.SugarLogger.Debugf("链接信息：\n%#v", cli.RemoteServerPID)
		case "to":
			if mt.Args.SubCommands[0] == "local" {
				cli.RemoteServerPID = cli.LocalServerPID
			} else {
				cli.RemoteServerPID = actor.NewPID(mt.Args.SubCommands[0], mt.Args.SubCommands[1])
			}
		}
		return
	case "transfer":
		go cli.transfer(mt.Command)
		return
	default:
		req = NewRequest(cli.LocalServerPID, mt)

		util.SugarLogger.Debugf("====================================================")
		util.SugarLogger.Debugf("任务内容：%s", req.ByteMsg)

		rstChan := make(chan []byte, 1)
		tasks.ResultCache.Add(req.UUID, rstChan)
	}

	// 收到isLocal的标志时, 在本地运行插件
	//if isLocal {
	//	var reverseFlag bool
	//	mt.Task.Args.Flags, reverseFlag = utils.IfInFlagThenPop("reverse", mt.Flags)
	//	pm := &tasks.PluginManager{
	//		PluginDir: pmt.PluginDir,
	//		PlugInfos: make(map[string]tasks.PlugInfo),
	//	}
	//	task := tasks.Task{PluginName: mt.PluginName, Cmd: mt.Cmd, SubCmd: mt.SubCmd,
	//		Flags: mt.Flags, Args: mt.Args}
	//
	//	//r, err := tasks.Exec(&task, pm)
	//	//utils.CheckErrorPanic(err)
	//	var (
	//		reply    *pb.CommonCmdReply
	//		nextTask *tasks.Task
	//		err      error
	//	)
	//	for {
	//		reply, nextTask, err = tasks.Exec(&task, pm, "", "")
	//		if nextTask == nil {
	//			break
	//		}
	//		task = *nextTask
	//		utils.CheckErrorPanic(err)
	//	}
	//
	//	err = CommonCmdOutputter(reply, reverseFlag)
	//	utils.CheckErrorPanic(err)
	//	return
	//}

	util.SugarLogger.Debugf("发送任务！")
	cli.RootContext.Send(cli.RemoteServerPID, req)
	util.SugarLogger.Debugf("发送任务完成！")

	if isRT {
		b := tasks.GetResult(req.UUID)
		util.SugarLogger.Debugf("任务返回值：\n%s", b)
	}

	//switch {
	//case mt.Type == pb.CommonCmdRequest_COMMON_CMD:
	//	var (
	//		reverseFlag bool
	//		//asyncFlag   bool
	//	)
	//	mt.Flags, reverseFlag = utils.IfInFlagThenPop("reverse", mt.Flags)
	//	//mt.Flags, asyncFlag = utils.IfInFlagThenPop("async", mt.Flags)
	//	//if asyncFlag {
	//	//	mt.Type = pb.CommonCmdRequest_ASYNC_TASK
	//	//}
	//	r, err := pmt.CC.CommonCmd(pmt.Ctx, &mt)
	//	utils.CheckErrorPanic(err)
	//	// 对help命令做特殊处理, 加载其中的智能提示, 并返回帮助信息
	//	if mt.Cmd == "help" {
	//		resEs := ParseYAML([]byte(r.ResultMsg))
	//		r.ResultMsg = resEs.Suggest.Description
	//		for k, v := range resEs.ESs {
	//			es.ESs[k] = v
	//		}
	//	}
	//
	//	err = CommonCmdOutputter(r, reverseFlag)
	//	utils.CheckErrorPanic(err)
	//case mt.Type == pb.CommonCmdRequest_FILE_TRANSFER && mt.PluginName == "upload":
	//	localPath := mt.Cmd
	//	filePath, fileName := filepath.Split(mt.Cmd)
	//	filePath = uploadPath(mt.Flags)
	//	applyFi, err := pmt.CC.ApplyTransfer(pmt.Ctx,
	//		&pb.TransferInfo{
	//			Type:       pb.TransferInfo_Upload,
	//			State:      pb.TransferInfo_Apply,
	//			FileName:   fileName,
	//			FilePath:   filePath,
	//			TransferId: 0,
	//		})
	//	utils.CheckErrorPanic(err)
	//	uploadStream, err := pmt.CC.Upload(pmt.Ctx)
	//	utils.CheckErrorPanic(err)
	//	// 创建一个1M的buf
	//	defer uploadStream.CloseSend()
	//	file, err := os.Open(localPath)
	//	utils.CheckErrorPanic(err)
	//	//stat, err := file.Stat()
	//	//utils.CheckErrorPanic(err)
	//	//fmt.Printf("文件大小: %d\n", stat.Size())
	//	buf := make([]byte, 1<<20)
	//	writing := true
	//	for writing {
	//		//fmt.Printf("读取文件 %s ...\n", fileName)
	//		n, err := file.Read(buf[:])
	//		if err != nil {
	//			if err == io.EOF {
	//				fmt.Println("文件已发送")
	//				writing = false
	//				err = nil
	//				continue
	//			}
	//			utils.CheckErrorPanic(err)
	//		}
	//		err = uploadStream.Send(&pb.Chunks{
	//			TransferId: applyFi.TransferId,
	//			Size:       int64(n),
	//			Content:    buf[:n],
	//		})
	//	}
	//	recvFi, err := uploadStream.CloseAndRecv()
	//	utils.CheckErrorPanic(err)
	//	checkUploadResult(recvFi, localPath, fileName)
	//}
	return
}

func (cli *Client) livePrefix() (string, bool) {
	return fmt.Sprintf("[%s] %s >>> ", cli.RemoteServerPID.Address, cli.Preload), true
}

func (cli *Client) Prepare() {
	cli.Prompt = prompt.New(
		cli.executor,
		completer,
		prompt.OptionPrefix(fmt.Sprintf("[%s-%s] %s >>> ", cli.RemoteServerPID.String(), cli.RemoteServerPID.Id, cli.Preload)),
		prompt.OptionLivePrefix(cli.livePrefix),
		prompt.OptionTitle("meepo"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
}

// 将命令行解析为 Request 格式
func Parse(line string) (*tasks.MeepoTask, bool, error) {
	//同步标志，所有的任务默认是异步的，加上标志的会转为同步任务
	isRT := false

	//分割命令
	fields, err := shlex.Split(line)
	if err != nil {
		util.SugarLogger.Debugf(err.Error())
		return nil, false, err
	}

	//空命令直接返回
	if len(fields) < 1 {
		return nil, false, nil
	}

	//执行本地命令
	if strings.ToLower(fields[0]) == "rt" {
		isRT = true
		if len(fields) < 2 {
			return nil, isRT, nil
		}
		fields = fields[1:]
	}

	//区分lua、lua脚本、task命令
	var mt *tasks.MeepoTask
	switch fields[0] {
	//case "lua":
	//	mt = tasks.NewMeepoTask("lua", fields[1], nil)
	//case "lua-script":
	//	mt = tasks.NewMeepoTask("lua-script", fields[1], nil)
	default:
		mt, err = parserTask(fields)
	}
	if err != nil {
		util.SugarLogger.Debugf(err.Error())
		return nil, isRT, err
	}
	return mt, isRT, nil
}

func NewRequest(sender *actor.PID, mt *tasks.MeepoTask) *pb.Request {
	req := NewRequestWithoutSender(mt)
	req.Sender = sender
	return req
}

func NewRequestWithoutSender(mt *tasks.MeepoTask) *pb.Request {
	b, err := mt.EncodeTask()
	if err != nil {
		util.SugarLogger.Errorf("task 编码错误：%s", err.Error())
	}

	return &pb.Request{
		UUID:    uuid.NewString(),
		Sender:  nil,
		MsgType: pb.MEEPO_TASK,
		ByteMsg: b,
	}
}

func parserTask(fields []string) (*tasks.MeepoTask, error) {
	var (
		//确定pluginName
		pluginName string
		command    string
		keyBuf     = ""
		valBuf     = make([]string, 0)
		args       = &tasks.Args{}
		cmdFields  []string
		argFields  []string
	)

	//分离前面的pluginName、command、subCmd和后面的args
	for idx, field := range fields {
		if strings.HasPrefix(field, "-") {
			cmdFields = fields[:idx]
			if len(fields) > idx {
				argFields = fields[idx:]
			}
			break
		}
	}

	if len(cmdFields) < 1 {
		cmdFields = fields
	}

	//解析pluginName、command、subCommand
	switch len(cmdFields) {
	case 0:
		return nil, errors.New("plugin name 不能为空")
	case 1:
		pluginName = cmdFields[0]
	case 2:
		pluginName = cmdFields[0]
		command = cmdFields[1]
	default:
		pluginName = cmdFields[0]
		command = cmdFields[1]
		args.SubCommands = cmdFields[2:]
	}

	//解析args
	for _, field := range argFields {
		if strings.HasPrefix(field, "-") {
			//归档Buf
			if keyBuf != "" {
				if len(valBuf) < 1 {
					args.Flags = append(args.Flags, keyBuf)
					keyBuf = ""
				} else {
					args.KVs[keyBuf] = valBuf
					keyBuf = ""
					valBuf = make([]string, 0)
				}
			}

			field = strings.TrimPrefix(field, "-")
			field = strings.TrimPrefix(field, "-")
			keyBuf = field
		} else {
			if keyBuf != "" {
				valBuf = append(valBuf, field)
			}
		}

	}

	//归档Buf
	if keyBuf != "" {
		if len(valBuf) < 1 {
			args.Flags = append(args.Flags, keyBuf)
		} else {
			args.KVs[keyBuf] = valBuf
		}
	}

	return tasks.NewMeepoTask(pluginName, command, args), nil
}
