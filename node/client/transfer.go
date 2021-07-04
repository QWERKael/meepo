package client

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"io"
	myactor "meepo/actor"
	"meepo/pb"
	"meepo/util"
	"os"
	"path/filepath"
	"time"
	"utility-go/path"
)


func (cli *Client) transfer2(filePath string) {
	_, fileName := filepath.Split(filePath)
	fileUUID := uuid.NewString()

	md5, err := path.SumMd5FromFile(filePath)
	if err != nil {
		util.SugarLogger.Errorf("无法获取文件的【%s】的MD5码：%s！", filePath, err.Error())
		return
	}

	trans := &pb.Transfer{
		FileUUID: fileUUID,
		Sender:   cli.LocalServerPID,
		State:    pb.Req,
		Size_:    int64(0),
		Context:  []byte(fileName+"|"+md5),
	}
	cli.RootContext.Send(cli.RemoteServerPID, trans)

	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		util.SugarLogger.Errorf("无法打开指定文件【%s】：%s！", filePath, err.Error())
		cli.sendClose(fileUUID)
		return
	}
	buf := make([]byte, 1<<20)
	writing := true
	for writing {
		state := pb.Data
		n, err := file.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				util.SugarLogger.Infof("文件【%s】已发送完成", filePath)
				writing = false
				state = pb.Done
				err = nil
				continue
			} else {
				util.SugarLogger.Errorf("无法读取文件【%s】：%s！", filePath, err.Error())
			}
		}
		trans.State = state
		trans.Context = buf
		trans.Size_ = int64(n)
		cli.RootContext.Send(cli.RemoteServerPID, trans)
	}
}

func (cli *Client) transfer(filePath string) {
	_, fileName := filepath.Split(filePath)
	fileUUID := uuid.NewString()

	md5, err := path.SumMd5FromFile(filePath)
	if err != nil {
		util.SugarLogger.Errorf("无法获取文件的【%s】的MD5码：%s！", filePath, err.Error())
		return
	}

	trans := &pb.Transfer{
		FileUUID: fileUUID,
		Sender:   cli.LocalServerPID,
		State:    pb.Req,
		Size_:    int64(0),
		Context:  []byte(fileName+"|"+md5),
	}
	var rst interface{}
	rst, err = cli.RootContext.RequestFuture(cli.RemoteServerPID, trans, 5*time.Second).Result()
	if err != nil {
		transferErrorProcess(err)
		return
	}
	util.SugarLogger.Debugf("文件传输【%s】接收到返回确认！", filePath)
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		util.SugarLogger.Errorf("无法打开指定文件【%s】：%s！", filePath, err.Error())
		cli.sendClose(fileUUID)
		return
	}
	buf := make([]byte, 1<<20)
	writing := true
	for writing {
		state := pb.Data
		n, err := file.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				util.SugarLogger.Infof("文件【%s】已发送完成", filePath)
				writing = false
				state = pb.Done
				err = nil
				continue
			} else {
				util.SugarLogger.Errorf("无法读取文件【%s】：%s！", filePath, err.Error())
			}
		}

		if r, ok := rst.(pb.Transfer); ok {
			switch r.State {
			case pb.OK:
				trans := &pb.Transfer{
					FileUUID: fileUUID,
					Sender:   cli.LocalServerPID,
					State:    state,
					Size_:    int64(n),
					Context:  buf[:n],
				}
				//rst, err = cli.RootContext.RequestFuture(cli.RemoteServerPID, trans, 5*time.Second).Result()
				rst, err = myactor.System.Root.RequestFuture(cli.RemoteServerPID, trans, 5*time.Second).Result()

				util.SugarLogger.Debugf("!!!%#v", myactor.System.Root.Sender().String())
				if err != nil {
					transferErrorProcess(err)
					cli.sendClose(fileUUID)
					return
				}
			case pb.Refuse:
				util.SugarLogger.Errorf("接收端拒绝接受文件【%s】！", filePath)
				cli.sendClose(fileUUID)
			case pb.Fail:
				util.SugarLogger.Errorf("接收端接受文件【%s】出错！", filePath)
				cli.sendClose(fileUUID)
			default:
				util.SugarLogger.Errorf("传输文件【%s】时发生未知错误！")
			}
		}
	}
	err = file.Close()
	if err != nil {
		util.SugarLogger.Errorf("无法关闭文件【%s】：%s！", filePath, err.Error())
	}

}

func transferErrorProcess(err error) {
	switch err {
	case actor.ErrTimeout:
		util.SugarLogger.Errorf("传输请求超时！")
	case actor.ErrDeadLetter:
		util.SugarLogger.Errorf("请求的PID错误！")
	default:
		util.SugarLogger.Errorf("未知的文件传输错误！")
	}
	return
}

func (cli *Client) sendClose(fileUUID string) {
	trans := &pb.Transfer{
		FileUUID: fileUUID,
		Sender:   cli.LocalServerPID,
		State:    pb.Close,
		Size_:    int64(0),
		Context:  nil,
	}
	cli.RootContext.Send(cli.RemoteServerPID, trans)
}
