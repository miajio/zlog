package zlog_test

import (
	"testing"

	"github.com/miajio/zlog"
)

func TestLoggerFile(t *testing.T) {
	l := zlog.Logger{
		Path:       "./log",
		MaxSize:    256,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   false,
	}
	lv := zlog.LogMap{
		"debug": zlog.DebufLevel,
		"info":  zlog.InfoLevel,
		"error": zlog.ErrorLevel,
	}

	l.Generate(lv)
	l.Log.Info("hello")
}
