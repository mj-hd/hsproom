package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	Level_Debug = iota
	Level_Info
	Level_Error
	Level_Fatal
)

type ErrorDetails struct {
	Message    string
	CallerFile string
	CallerLine int
	Level      int
}

var LogFile string
var DisplayLog bool
var LogLevel int

func PromulgateInfoStr(w io.Writer, message string) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    message,
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Info,
	})
}

func PromulgateInfo(w io.Writer, err error) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    err.Error(),
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Info,
	})
}

func PromulgateDebugStr(w io.Writer, message string) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    message,
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Debug,
	})
}

func PromulgateDebug(w io.Writer, err error) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    err.Error(),
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Debug,
	})
}

func PromulgateErrorStr(w io.Writer, message string) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    message,
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Error,
	})
}

func PromulgateError(w io.Writer, err error) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    err.Error(),
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Error,
	})
}

func PromulgateFatalStr(w io.Writer, message string) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    message,
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Fatal,
	})
}

func PromulgateFatal(w io.Writer, err error) {
	_, file, line, _ := runtime.Caller(1)
	Promulgate(w, ErrorDetails{
		Message:    err.Error(),
		CallerFile: file,
		CallerLine: line,
		Level:      Level_Fatal,
	})
}

func Promulgate(w io.Writer, details ErrorDetails) {

	if details.Level < LogLevel {
		return
	}

	var level string

	switch details.Level {
	case Level_Info:
		level = "Info"
	case Level_Debug:
		level = "Debug"
	case Level_Error:
		level = "Error"
	case Level_Fatal:
		level = "Fatal"
	}

	log := fmt.Sprintf("[%s]%s(%d)[%s]: %s\n", time.Now().Format("Jan _2 15:04:05"), details.CallerFile, details.CallerLine, level, details.Message)

	if DisplayLog {
		w.Write([]byte(log))
	}

	file, err := os.OpenFile(LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(log)

	return
}
