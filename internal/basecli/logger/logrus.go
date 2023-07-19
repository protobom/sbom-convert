package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// LogrusConfig contains all configurable values for the Logrus logger
type LogrusConfig struct {
	EnableConsole bool
	EnableFile    bool
	Structured    bool
	Level         logrus.Level
	FileLocation  string
}

// NewLogrusLogger creates a new LogrusLogger with the given configuration
func NewLogrusLogger(cfg LogrusConfig) (Logger, error) {
	appLogger := logrus.New()
	var output io.Writer
	switch {
	case cfg.EnableConsole && cfg.EnableFile:
		logFile, err := os.OpenFile(cfg.FileLocation, os.O_WRONLY|os.O_CREATE, createPerm)
		if err != nil {
			return nil, fmt.Errorf("unable to setup log file: %w", err)
		}
		output = io.MultiWriter(os.Stderr, logFile)
	case cfg.EnableConsole:
		output = os.Stderr
	case cfg.EnableFile:
		logFile, err := os.OpenFile(cfg.FileLocation, os.O_WRONLY|os.O_CREATE, createPerm)
		if err != nil {
			return nil, fmt.Errorf("unable to setup log file: %w", err)
		}
		output = logFile
	default:
		output = io.Discard
	}

	appLogger.SetOutput(output)
	appLogger.SetLevel(cfg.Level)
	if cfg.Structured {
		appLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			PrettyPrint:       false,
		})
	} else {
		appLogger.SetFormatter(&prefixed.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			ForceFormatting: true,
			FullTimestamp:   true,
		})
	}

	return appLogger, nil
}
