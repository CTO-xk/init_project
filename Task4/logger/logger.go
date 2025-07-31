package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger(env, logPath string) error {
	var config zap.Config

	// 根据环境选择不同的日志配置
	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 创建多个写入器：控制台 + 文件
	fileWriter := getLogWriter(logPath)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 定义不同级别的输出目标
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(getEncoder(env), fileWriter, highPriority),
		zapcore.NewCore(getEncoder(env), consoleWriter, lowPriority),
	)

	// 构建日志实例
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(Logger) // 替换全局日志实例

	return nil
}

// getEncoder 根据环境返回不同的编码器
func getEncoder(env string) zapcore.Encoder {
	if env == "production" {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// getLogWriter 返回带滚动功能的文件写入器
func getLogWriter(logPath string) zapcore.WriteSyncer {
	if logPath == "" {
		logPath = "./logs/sepolia-ops.log"
	}

	// 创建日志目录（如果不存在）
	if err := os.MkdirAll("./logs", 0755); err != nil {
		zap.L().Error("创建日志目录失败", zap.Error(err))
		return nil
	}

	// 使用 lumberjack 实现日志滚动
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,   // 单个文件最大 10MB
		MaxBackups: 30,   // 最多保留 30 个备份
		MaxAge:     7,    // 最多保留 7 天
		Compress:   true, // 启用压缩
	}

	return zapcore.AddSync(lumberjackLogger)
}

// Sync 同步日志缓冲区（程序退出前调用）
func Sync() {
	_ = Logger.Sync()
}
