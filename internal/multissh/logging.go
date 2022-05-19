package multissh

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogger : log formatting bits
func InitLogger(debug bool, format string) {
    // Set log level
	logLevel := log.InfoLevel
	if debug {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

    // Set format
    var formatter log.Formatter
    switch format {
    case "human":
        formatter = &log.TextFormatter{}
    case "json":
        formatter = &log.JSONFormatter{}
    }
	log.SetFormatter(formatter)

    // Set output
	log.SetOutput(os.Stdout)
}

func logger(stage string) *log.Entry {
    return log.WithFields(log.Fields{
        "stage": stage,
    })
}

