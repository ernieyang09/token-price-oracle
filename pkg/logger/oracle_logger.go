package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)


var log *logrus.Logger

func InitLogger(filePath string) {
	log = logrus.New()

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    20, // megabytes
		MaxBackups: 5,
		MaxAge:     28,    // days
		Compress:   false, // disabled by default
	}

	log.SetFormatter(&logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime: "timestamp",
            logrus.FieldKeyMsg:  "message",
            logrus.FieldKeyFunc: "func",
            logrus.FieldKeyFile: "file",
        },
    })

	multiWriter := io.MultiWriter(lumberjackLogger, os.Stdout)
    log.SetOutput(multiWriter)


}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}