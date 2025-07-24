package infra

import (
	"log"
	"os"
)

var infoLog = log.New(os.Stdout, "[yact] ", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "[yact] ", log.Lshortfile)

func Info(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Error(format string, v ...any) {
	errorLog.Printf(format, v...)
}
