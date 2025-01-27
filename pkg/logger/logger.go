package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New(os.Stdout, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	Logger.Println("INFO: " + message)
}

func Error(message string) {
	Logger.Println("ERROR: " + message)
}

func Warn(message string) {
	Logger.Println("WARN: " + message)
}
