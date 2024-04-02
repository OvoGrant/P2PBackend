package main

import (
	"log"
	"os"
)

var (
	ConnectionLogger *log.Logger
	ErrorLogger      *log.Logger
	IndexingLogger   *log.Logger
	DeletionLogger   *log.Logger
)

// initLoggers creates loggers that write to a text file logs.txt
func initLoggers() {

	file, err := os.Create("logs.txt")

	if err != nil {
		log.Fatalf(err.Error())
	}

	ConnectionLogger = log.New(file, "CONNECTION: ", log.Ldate|log.Ltime|log.Lshortfile)
	IndexingLogger = log.New(file, "INDEXING: ", log.Ldate|log.Ltime|log.Lshortfile)
	DeletionLogger = log.New(file, "DELETION: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
