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

func InitLogger(config config.LoggerConf) error {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{})
	switch config.Level {
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
		config.Path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(config.Path),
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)
	if err != nil {
		return err
	}

	if config.Stdout {
		writer = io.MultiWriter(os.Stdout, fileWriter)
	} else {
		writer = fileWriter
	}
	Logger.SetOutput(writer)
	return nil
}
