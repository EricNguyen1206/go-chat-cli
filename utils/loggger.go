package utils

import (
	"log"

	"github.com/fatih/color"
)

var (
	info  = color.New(color.FgGreen).SprintFunc()
	warn  = color.New(color.FgYellow).SprintFunc()
	err   = color.New(color.FgRed).SprintFunc()
	debug = color.New(color.FgCyan).SprintFunc()
)

func Info(v ...any) {
	log.Println(info("INFO:"), v)
}

func Warn(v ...any) {
	log.Println(warn("WARN:"), v)
}

func Error(v ...any) {
	log.Println(err("ERROR:"), v)
}

func Debug(v ...any) {
	log.Println(debug("DEBUG:"), v)
}
