package logger

import "log"

func Info(msg string) {
	log.Println("[INFO]:", msg)
}

func Error(msg string) {
	log.Println("[ERROR]:", msg)
}

func Panic(msg string) {
	log.Panicln("[PANIC]:", msg)
}
