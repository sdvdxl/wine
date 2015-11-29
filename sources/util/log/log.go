package log

import (
	"github.com/sdvdxl/log4go"
)

var (
	Logger log4go.Logger
)

func init() {
	//初始化日志
	Logger = log4go.NewLogger()
	Logger.AddFilter("fine", log4go.FINE, log4go.NewConsoleLogWriter())
	Logger.AddFilter("error", log4go.ERROR, log4go.NewFileLogWriter("error.log", false))
	Logger.AddFilter("info", log4go.INFO, log4go.NewFileLogWriter("info.log", false))
	Logger.Info("log inited")
}

func Close() {
	Logger.Close()
}
