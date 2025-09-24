package logger

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

// New 创建新的日志实例
func New(module string) *slog.Logger {
	// 开发环境使用文本格式，生产环境可改为JSON格式
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, ok := a.Value.Any().(*slog.Source)
				if ok {
					// 简化文件名
					source.File = filepath.Base(source.File)
					// 隐藏日志包自身的调用栈
					if source.File == "logger.go" {
						return slog.Attr{}
					}
				}
			}
			return a
		},
	})

	// 添加模块名
	return slog.New(handler).With("module", module)
}

// Fatal 致命错误并退出
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	log.Fatal(msg)
}
