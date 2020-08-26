package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

const _DATE_FORMAT = "2006-01-02 15:04:05.000"

func debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("DEBUG", msg)
}

func info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("INFO", msg)
}

func errs(v ...interface{}) {
	msg := fmt.Sprint(v...)
	_log("ERROR", msg)
}

func errsf(formate string, v ...interface{}) {
	msg := fmt.Sprintf(formate, v...)
	_log("ERROR", msg)
}

func _log(level string, msg string) {
	now := time.Now().Format(_DATE_FORMAT)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	log.Print(fmt.Sprintf("[%s] [Thread-0] [%s] [%s:%d] [TID:N/A] - %s\n", now, level, file, line, msg))
}
