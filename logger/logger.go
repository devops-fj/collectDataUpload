// logger/logger.go
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

// InitLogger 初始化日志
func InitLogger(level string, output string) {
	log = logrus.New()

	// 设置日志级别
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// 设置日志输出
	switch output {
	case "console":
		log.SetOutput(os.Stdout)
	case "file":
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	default:
		log.SetOutput(os.Stdout)
	}
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return log
}
