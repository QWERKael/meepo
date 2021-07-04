package common

type TransferState uint8

const (
	Unknown TransferState = 0 // 未知状态
	Request TransferState = 1 // 发送者向接收者请求传输文件
	OK      TransferState = 2 // 接收者同意传输文件，或者接收者成功接收到文件块
	Refuse  TransferState = 3 // 接收者拒绝接收文件，可能是因为文件有重复
	Data    TransferState = 4 // 发送者传输文件块
	Done    TransferState = 5 // 发送者发送的最后一个文件块
)
