# zlog
Uber based zaplog carries out simple and crude encapsulation, allowing users to quickly develop log modules

# Install
go get -u github.com/miajio/zlog

```
package zlog_test

import (
	"testing"

	"github.com/miajio/zlog"
)

func TestLoggerFile(t *testing.T) {
    // set log param
	l := zlog.Logger{
		Path:       "./log",
		MaxSize:    256,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   false,
	}
    // set log level func
	lv := zlog.LogMap{
		"debug": zlog.DebufLevel,
		"info":  zlog.InfoLevel,
		"error": zlog.ErrorLevel,
	}
    // init logger
	l.Generate(lv)
    // info
	l.Log.Info("hello")
}

```