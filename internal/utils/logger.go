package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(logDir string, logLevel string) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", logLevel)
	}
	Log.SetLevel(level)

	// 设置日志格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 设置日志输出
	logPath := filepath.Join(logDir, "eth-trading-system.log")
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize log file: %v", err)
	}

	// 同时将日志输出到控制台和文件
	Log.SetOutput(io.MultiWriter(os.Stdout, writer))

	return nil
}
