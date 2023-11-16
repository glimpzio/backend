package misc

import (
	"fmt"
	"io"
	"log"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Create a new logger
func NewLogger(name string, out io.Writer) *Logger {
	return &Logger{InfoLog: log.New(out, fmt.Sprintf("[%s] INFO: ", name), log.Ldate|log.Ltime), ErrorLog: log.New(out, fmt.Sprintf("[%s] ERROR: ", name), log.Ldate|log.Ltime)}
}
