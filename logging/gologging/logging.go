/**
 * @Author: huangw1
 * @Date: 2019/7/10 20:12
 */

package gologging

import (
	gologging "github.com/op/go-logging"
	"github.com/huangw1/gateway/logging"
	"io"
)

const module = "gateway"

func NewLogger(level string, out io.Writer, prefix string) logging.Logger {
	log := gologging.MustGetLogger(module)
	logBackend := gologging.NewLogBackend(out, prefix, 0)
	format := gologging.MustStringFormatter(
		` %{time:2006/01/02 - 15:04:05.000} %{color}▶ %{level:.6s}%{color:reset} %{message}`,
	)
	backendFormatter := gologging.NewBackendFormatter(logBackend, format)
	backendLeveled := gologging.AddModuleLevel(backendFormatter)
	logLevel, _ := gologging.LogLevel(level)
	backendLeveled.SetLevel(logLevel, module)
	gologging.SetBackend(backendLeveled)
	return Logger{log}
}

type Logger struct {
	Logger *gologging.Logger
}

func (l Logger) Debug(v ...interface{}) {
	l.Logger.Debug(v...)
}

func (l Logger) Info(v ...interface{}) {
	l.Logger.Info(v...)
}

func (l Logger) Warning(v ...interface{}) {
	l.Logger.Warning(v...)
}

func (l Logger) Error(v ...interface{}) {
	l.Logger.Error(v...)
}

func (l Logger) Critical(v ...interface{}) {
	l.Logger.Critical(v...)
}

func (l Logger) Fatal(v ...interface{}) {
	l.Logger.Fatal(v...)
}

