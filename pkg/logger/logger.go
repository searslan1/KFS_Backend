package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger nesnesi
var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}

func Info(msg string) {
	Log.Info(msg)
}

func Error(msg string) {
	Log.Error(msg)
}
