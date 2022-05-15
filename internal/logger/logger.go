package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/saying-yan/embedded_system_course_project_backend/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var Logger *logrus.Logger

func InitLoggerWithConf(config config.LoggerConf) error {
	return InitLogger(config.Level, config.Path, config.Stdout)
}

func InitLogger(level, path string, stdout bool) error {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{})
	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	var writer io.Writer
	fileWriter, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)
	if err != nil {
		return err
	}

	if stdout {
		writer = io.MultiWriter(os.Stdout, fileWriter)
	} else {
		writer = fileWriter
	}
	Logger.SetOutput(writer)
	return nil
}
