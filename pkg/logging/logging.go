package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/supporttools/prometheus-tunnel/pkg/config"
)

var log = SetupLogging(config.CFG.Debug) // Initialize the logger

// LogFile returns a log entry with the filename and line number
func LogFile() *logrus.Entry {
	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		log.Panic("Unable to get caller information")
	}
	filename = filepath.Base(filename)

	// Check if the logger is in debug mode
	if config.CFG.Debug {
		fmt.Println("Debug logging enabled")
		return log.WithFields(logrus.Fields{"filename": filename, "line": line})
	}

	// If not in debug mode, return a log entry without the filename
	return log.WithField("line", line)
}

// SetupLogging initializes the logger with the specified log level
func SetupLogging(debug bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)

	customFormatter := new(logrus.TextFormatter)
	// To remove timestamp, set TimestampFormat to an empty string
	customFormatter.TimestampFormat = ""
	customFormatter.FullTimestamp = false // Also disable the full timestamp to ensure no timestamp is printed
	logger.SetFormatter(customFormatter)

	logger.SetOutput(os.Stderr)

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// GetRelativePath returns the relative path of a file from the working directory
func GetRelativePath(filePath string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Error("Failed to get working directory")
		return filePath // Returning the original path as a fallback
	}
	relPath, err := filepath.Rel(wd, filePath)
	if err != nil {
		log.WithError(err).WithField("filePath", filePath).Error("Failed to get relative path")
		return filePath // Returning the original path as a fallback
	}
	return relPath
}
