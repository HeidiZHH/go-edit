package main

import (
	"log"
	"os"

	"gihub.com/heidizhh/go-edit/internal/editor"
)

func main() {
	logFile, err := os.OpenFile("edit.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logFile)
	e, err := editor.New()
	if err != nil {
		log.Fatalf("could not create editor: %v", err)
	}
	defer e.Close()
	if err := e.Run(); err != nil {
		log.Fatalf("editor aborted: %v", err)
	}
}
