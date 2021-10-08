package logger

import (
	"log"
	"os"
)

var (
	// WarningLogger contains template for warning logs
	WarningLogger *log.Logger
	// InfoLogger contains template for info logs
	InfoLogger *log.Logger
	// ErrorLogger contains template for error logs
	ErrorLogger *log.Logger
)

// InitLogger inits loggers
func InitLogger() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
