package common

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
	"os"
)

var Mee = &Meepo{}

type Meepo struct {
	L              *lua.LState
	RootContext    *actor.RootContext
	LocalServerPID *actor.PID
}

type TransferHeader struct {
	TransferState TransferState
	Name          string
	UUID          string
}

type FileBuilder struct {
	FileUUID string
	Sender   *actor.PID
	Name     string
	MD5      string
	File     *os.File
	Channel  chan []byte
	Done     bool
}
