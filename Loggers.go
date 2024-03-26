package main

import (
	"log"
	"os"
)

var (
	ConnectionLogger *log.Logger
	IndexingLogger   *log.Logger
	DeletionLogger   *log.Logger
)

func initLoggers() {
	ConnectionLogger = log.New(os.Stdout, "CONNECTION: ", log.Ldate|log.Ltime|log.Lshortfile)
	IndexingLogger = log.New(os.Stdout, "INDEXING: ", log.Ldate|log.Ltime|log.Lshortfile)
	DeletionLogger = log.New(os.Stdout, "DELETION: ", log.Ldate|log.Ltime|log.Lshortfile)
}
