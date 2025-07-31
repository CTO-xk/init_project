package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CTO-xk/init_project/Task4/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志实例
var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// Init 初始化日志系统
func Init(cfg config.LogConfig) error {
	// 验证配置
	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("日志配置无效: %w", err)
	}

	// 解析日志级别
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}

	// 创建编码器
	encoder := newEncoder(cfg.Env)

	// 创建输出目标
	writers := buildWriters(cfg)

	if len(writers) == 0 {
		return fmt.Errorf("至少需要配置一个日志输出目标")
	}

	// 创建日志核心
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)

	// 构建日志实例
	Logger = zap.New(
		core,
		zap.AddCaller(),                       // 记录调用位置
		zap.AddCallerSkip(1),                  // 跳过当前函数帧
		zap.AddStacktrace(zapcore.ErrorLevel), // 错误级别记录堆栈
	)
	Sugar = Logger.Sugar()

	Logger.Info("日志系统初始化完成",
		zap.String("env", cfg.Env),
		zap.String("level", cfg.Level),
	)
	return nil
}

// 验证日志配置
func validateConfig(cfg config.LogConfig) error {
	if cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}
	}
	if cfg.MaxSize <= 0 {
		return fmt.Errorf("max_size必须大于0")
	}
	return nil
}

// 解析日志级别
func parseLevel(levelStr string) (zapcore.Level, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(levelStr))); err != nil {
		return level, fmt.Errorf("无效的日志级别: %s, 可选值: debug/info/warn/error/fatal", levelStr)
	}
	return level, nil
}

// 创建编码器
func newEncoder(env string) zapcore.Encoder {
	if env == "production" {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	// 开发环境使用带颜色的控制台编码器
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 构建输出目标
func buildWriters(cfg config.LogConfig) []zapcore.WriteSyncer {
	var writers []zapcore.WriteSyncer

	// 开发环境默认输出到控制台
	if cfg.Env == "development" {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 添加文件输出
	if cfg.FilePath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		})
		writers = append(writers, fileWriter)
	}

	return writers
}

// Sync 同步日志到磁盘
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

// 封装常用日志方法
func Debug(msg string, fields ...zap.Field) { Logger.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { Logger.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Logger.Warn(msg, fields...) }
func Error(msg string, err error, fields ...zap.Field) {
	Logger.Error(msg, append(fields, zap.Error(err))...)
}
func Fatal(msg string, err error, fields ...zap.Field) {
	Logger.Fatal(msg, append(fields, zap.Error(err))...)
}
