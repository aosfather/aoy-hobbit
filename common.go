package main

import (
	"fmt"
	"github.com/aosfather/bingo_utils"
	"log"
	"runtime"
	"time"
)

const _DATE_FORMAT = "2006-01-02 15:04:05.000"

func init() {
	log.SetFlags(0)
}
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
	now := time.Now().Format(bingo_utils.FORMAT_DATETIME_LOG)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	var short string
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	log.Print(fmt.Sprintf("[%s] [%s] [%s:%d] - %s\n", now, level, short, line, msg))
}
