package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var (
	Http, Info, Debug, Error *log.Logger
)

func init() {
	writer, err := rotatelogs.New(
		"./logs/%y-%m-%d.log",
		rotatelogs.WithLinkName(""),
		rotatelogs.WithMaxAge(time.Duration(7*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(1*24)*time.Hour),
	)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	Http = log.New(writer, "[Http]\t", log.Ldate|log.Ltime)
	Info = log.New(writer, "[Info]\t", log.Ldate|log.Ltime)
	Debug = log.New(writer, "[Debug]\t", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(writer, "[Error]\t", log.Ldate|log.Ltime|log.Lshortfile)
}
