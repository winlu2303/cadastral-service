package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(level string) {
	flags := log.LstdFlags | log.Lshortfile

	Info = log.New(os.Stdout, "INFO: ", flags)
	Error = log.New(os.Stderr, "ERROR: ", flags)

	if level == "debug" {
		Info.SetFlags(flags | log.Lmicroseconds)
		Error.SetFlags(flags | log.Lmicroseconds)
	}
}
