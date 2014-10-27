package debug

import (
	"log"
)

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}
