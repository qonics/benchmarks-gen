package helper

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func logFile() *os.File {
	file, err := os.OpenFile("cache-manager-"+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Open log file failed: %v", err))
	}
	return file
}
func singleLogFile(title string) *os.File {
	file, err := os.OpenFile("cache-manager-"+title+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Open single %v log file failed: %v", title, err))
	}
	return file
}
func Warning(message string) {
	file := logFile()
	log.SetOutput(file)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println(message)
}
func Critical(message string, title string) {
	file := singleLogFile("critical")
	log.SetOutput(file)
	WarningLogger = log.New(file, title+": ", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println(message)
}
