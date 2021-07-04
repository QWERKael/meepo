package util

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdLog "log"
	"utility-go/log"
)

var SugarLogger *zap.SugaredLogger
var StdLogger *stdLog.Logger

func init() {
	fmt.Println("初始化模块: log.go")
	var err error
	SugarLogger, err = log.NewLogger(log.ConsoleEncoder, "", zapcore.DebugLevel)
	if err != nil {
		panic(err.Error())
	}
	SugarLogger.Debug("日志记录开始...")
}

func LogInit() {
	var err error
	SugarLogger, err = log.NewLogger(log.ConsoleEncoder, Config.LogPathAbs, zapcore.DebugLevel)
	if err != nil {
		panic(err.Error())
	}
	StdLogger, err = zap.NewStdLogAt(SugarLogger.Desugar(), zapcore.DebugLevel)
	if err != nil {
		panic(err.Error())
	}
	SugarLogger.Debug("日志记录开始...")
}
