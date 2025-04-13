package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	defaultLogger = log.New(os.Stdout, "LOG: ", log.LstdFlags)
	errorLogger   = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func Info(args ...any) {
	defaultLogger.Println(args...)
}

func Error(args ...any) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		errorLogger.Println(fmt.Sprintf("%s:%d", file, line), args)
	} else {
		errorLogger.Println(args...)
	}
}
