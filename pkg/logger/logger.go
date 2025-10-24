package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLogger  = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
)

func Info(msg string) {
	logWithLevel("INFO", msg, infoLogger)
}

func Error(msg string) {
	logWithLevel("ERROR", msg, errorLogger)
}

func logWithLevel(level, msg string, logger *log.Logger) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level, msg)
	logger.Println(formattedMsg)
}
