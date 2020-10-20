package zapconf

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger *zap.Logger
	debugOnce sync.Once
	infoOnce sync.Once
	warnOnce sync.Once
	errorOnce sync.Once
	debugInstance *zap.Logger
	infoInstance *zap.Logger
	warnInstance *zap.Logger
	errorInstance *zap.Logger
)

// logPath 日志文件路径
// logLevel 日志级别 debug/info/warn/error
// maxSize 单个文件大小,MB
// maxBackups 保存的文件个数
// maxAge 保存的天数， 没有的话不删除
// compress 压缩
// jsonFormat 是否输出为json格式
// shoowLine 显示代码行
// logInConsole 是否同时输出到控制台

func InitLogger(logPath string, logLevel string, maxSize, maxBackups, maxAge int, compress, jsonFormat, showLine, logInConsole bool) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // 最多保留300个备份
		Compress:   compress,   // 是否压缩 disabled by default
	}
	if maxAge > 0 {
		hook.MaxAge = maxAge // days
	}

	var syncer zapcore.WriteSyncer
	if logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	logger = zap.New(core)
	if showLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func GetDebugLog() *zap.Logger {
	debugOnce.Do(func() {
		debugInstance = InitLogger(LogPath, "debug", LogMaxSize, LogMaxBackups, 0, false, true, true, true)
	})
	return debugInstance
}

func GetInfoLog() *zap.Logger {
	infoOnce.Do(func() {
		infoInstance = InitLogger(LogPath, "info", LogMaxSize, LogMaxBackups, 0, false, true, true, true)
	})
	return infoInstance
}

func GetWarnLog() *zap.Logger {
	warnOnce.Do(func() {
		warnInstance = InitLogger(LogPath, "warn", LogMaxSize, LogMaxBackups, 0, false, true, true, true)
	})
	return warnInstance
}

func GetErrorLog() *zap.Logger {
	errorOnce.Do(func() {
		errorInstance = InitLogger(LogPath, "error", LogMaxSize, LogMaxBackups, 0, false, true, true, true)
	})
	return errorInstance
}