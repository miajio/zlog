package zlog

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level log level
type Level func(zapcore.Level) bool

// LogMap log level map
// map key is file name
// value is log print level func
type LogMap map[string]Level

// Logger
type Logger struct {
	Path       string             `json:"path" toml:"toml"`             // log file path
	MaxSize    int                `json:"maxSize" toml:"maxSize"`       // log file max size
	MaxBackups int                `json:"maxBackups" toml:"maxBackups"` // log file max backups
	MaxAge     int                `json:"maxAge" toml:"maxAge"`         // log file max save day
	Compress   bool               `json:"compress" toml:"compress"`     // log file whether to compress
	logMap     *LogMap            // log level map
	Log        *zap.SugaredLogger // zap log object
	mu         sync.Mutex         // logger init lock
}

// LoggerInterface
type LoggerInterface interface {
	Generate(LogMap) // Generate logger core
}

// default log level
var (
	DebufLevel = func(level zapcore.Level) bool {
		return level == zap.DebugLevel
	}

	InfoLevel = func(level zapcore.Level) bool {
		return level == zap.InfoLevel
	}

	ErrorLevel = func(level zapcore.Level) bool {
		return level == zap.ErrorLevel
	}

	_ LoggerInterface = (*Logger)(nil)
)

// Generate
func (log *Logger) Generate(logMap LogMap) {
	if log.Log == nil {
		log.mu.Lock()
		defer log.mu.Unlock()
		if log.Log == nil {
			encoderConfig := zapcore.EncoderConfig{
				TimeKey:       "time",
				LevelKey:      "level",
				NameKey:       "log",
				CallerKey:     "lineNum",
				MessageKey:    "msg",
				StacktraceKey: "stacktrace",
				LineEnding:    zapcore.DefaultLineEnding,
				EncodeLevel:   zapcore.LowercaseLevelEncoder,
				EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
					pae.AppendString(t.Format("[2006-01-02 15:04:05]"))
				},
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
				EncodeName:     zapcore.FullNameEncoder,
			}

			cores := make([]zapcore.Core, 0)

			for fileName := range logMap {
				if fileName == "" {
					continue
				}
				// logger write
				write := &lumberjack.Logger{
					Filename:   GetLogFilePath(log.Path, fileName),
					MaxSize:    log.MaxSize,
					MaxBackups: log.MaxBackups,
					MaxAge:     log.MaxAge,
					Compress:   log.Compress,
				}
				// the log level
				level := zap.LevelEnablerFunc(logMap[fileName])
				// log core
				core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(write)), level)
				cores = append(cores, core)
			}
			// append default info log level
			cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel))

			core := zapcore.NewTee(cores...)
			caller := zap.AddCaller()
			development := zap.Development()

			log.Log = zap.New(core, caller, development, zap.Fields()).Sugar()
		}
	}
}

// logFile logger file out path
func GetLogFilePath(filePath, fileName string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	fp := strings.Split(filePath, "/")
	realPath := make([]string, 0)
	for i := range fp {
		if fp[i] != "" {
			realPath = append(realPath, fp[i])
		}
	}

	filePath = strings.Join(realPath, "/")
	if filePath == "" {
		filePath = "."
	}

	if !strings.HasSuffix(fileName, ".log") {
		fileName = strings.ReplaceAll(fileName, ".", "_")
		fileName = fileName + ".log"
	}
	return filePath + "/" + fileName
}
