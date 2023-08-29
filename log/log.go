package log

import (
	"log"
	"os"
)

var (
	InfoLogger = log.New(os.Stdout, "[I] ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(os.Stdout, "[W] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLogger  = log.New(os.Stdout, "[E] ", log.Ldate|log.Ltime|log.Lshortfile)
)
