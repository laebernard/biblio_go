package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

func init() {
	// Créer les fichiers de log s'ils n'existent pas
	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.Mkdir("logs", 0755)
		infoFile, _ = os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		errorFile, _ = os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	infoLog = log.New(infoFile, "[INFO] ", log.LstdFlags)
	errorLog = log.New(errorFile, "[ERROR] ", log.LstdFlags)
}

func Info(message string, args ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	infoLog.Printf(msg, args...)
}

func Error(message string, args ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	errorLog.Printf(msg, args...)
}

func InfoWithContext(context string, message string) {
	Info("[%s] %s", context, message)
}

func ErrorWithContext(context string, message string) {
	Error("[%s] %s", context, message)
}
