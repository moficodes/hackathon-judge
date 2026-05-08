package logger

import (
	"log"
	"os"
)

// Init configures the standard library logger.
func Init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
