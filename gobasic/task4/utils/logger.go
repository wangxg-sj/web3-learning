package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() error {
	// 获取环境变量,默认为 development
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建不同级别的日志文件写入器
	infoWriter := getLogWriter(filepath.Join(logDir, "info.log"))
	errorWriter := getLogWriter(filepath.Join(logDir, "error.log"))
	debugWriter := getLogWriter(filepath.Join(logDir, "debug.log"))

	// 创建不同级别的 Core
	var cores []zapcore.Core

	if env == "development" {
		// 开发环境:控制台彩色输出 + 文件输出
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores,
			// Debug 级别及以上输出到控制台
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			// Debug 级别输出到 debug.log
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), debugWriter, zapcore.DebugLevel),
		)
	}

	// Info 级别输出到 info.log
	cores = append(cores,
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			infoWriter,
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl == zapcore.InfoLevel
			}),
		),
	)

	// Error 级别及以上输出到 error.log
	cores = append(cores,
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			errorWriter,
			zapcore.ErrorLevel,
		),
	)

	// 创建 Logger
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// getLogWriter 创建日志文件写入器(支持自动切割)
func getLogWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,  // 单个文件最大 100MB
		MaxBackups: 10,   // 最多保留 10 个备份
		MaxAge:     30,   // 保留 30 天
		Compress:   true, // 压缩旧日志
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Sync 刷新日志缓冲区
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Info 记录 Info 级别日志
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		// 跳过一层调用栈,显示真实的调用位置
		Logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
	}
}

// Error 记录 Error 级别日志
func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		// 跳过一层调用栈,显示真实的调用位置
		Logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
	}
}

// Debug 记录 Debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		// 跳过一层调用栈,显示真实的调用位置
		Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
	}
}

// Warn 记录 Warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		// 跳过一层调用栈,显示真实的调用位置
		Logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
	}
}

// Fatal 记录 Fatal 级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		// 跳过一层调用栈,显示真实的调用位置
		Logger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
	}
}
