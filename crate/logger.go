// Handles logging to a file in the Crate data directory
// Log line format is as follows (for regular expression parsing)
// 		%(level)s [%(jsontime)s]: %(message)s

package crate

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bbengfort/crate/crate/config"
)

var eventLogger *Logger // global var for the logger object

//=============================================================================

// LogLevel types
type LogLevel int

const (
	LevelDebug LogLevel = 1 + iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levels = [...]string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
}

func (level LogLevel) String() string {
	return levels[level-1]
}

func LevelFromString(level string) LogLevel {
	level = strings.ToUpper(level)
	level = strings.Trim(level, " ")

	switch level {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelInfo
	}
}

//=============================================================================

// Wraps the log.Logger to provide custom log handling
type Logger struct {
	Level   LogLevel    // The minimum log level to log at
	writer  *log.Logger // The logger object that handles logging
	logfile *os.File    // Handle to the open log file
}

//=============================================================================

// Initialize the Logger objects for logging to config location
func InitializeLoggers(level LogLevel) error {

	// Create a new logger
	eventLogger = new(Logger)
	eventLogger.Level = level

	// Open a handle to the log file
	path, err := config.CrateLoggingPath()
	if err != nil {
		return err
	}

	// Set the output path (opening the file and configuring the writer)
	return eventLogger.SetOutputPath(path)
}

// Close the loggers (and the open file handle) useful for defering close
func CloseLoggers() error {
	err := eventLogger.Close()
	eventLogger = nil
	return err
}

// Write a log message to the eventLogger at a certain log level
func Log(msg string, level LogLevel, args ...interface{}) {
	eventLogger.Log(msg, level, args...)
}

//=============================================================================

// Close the open handle to the log file and stop logging
func (logger *Logger) Close() error {
	return logger.logfile.Close()
}

// Set a new log output location on the Logger
func (logger *Logger) SetOutputPath(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	eventLogger.logfile = file
	eventLogger.writer = log.New(eventLogger.logfile, "", 0)

	return nil
}

// Write a log message to the logger with a certain log level
func (logger *Logger) Log(layout string, level LogLevel, args ...interface{}) {
	// Log line format is "%(level)s [%(jsontime)s]: %(message)s"

	if level >= logger.Level {
		msg := fmt.Sprintf(layout, args...)
		logger.writer.Printf("%-7s [%s]: %s\n", level, time.Now().Format(JSONLayout), msg)
	}
}

// Helper function to log at debug level
func (logger *Logger) Debug(msg string, args ...interface{}) {
	logger.Log(msg, LevelDebug, args...)
}

// Helper function to log at info level
func (logger *Logger) Info(msg string, args ...interface{}) {
	logger.Log(msg, LevelInfo, args...)
}

// Helper function to log at debug level
func (logger *Logger) Warn(msg string, args ...interface{}) {
	logger.Log(msg, LevelWarn, args...)
}

// Helper function to log at debug level
func (logger *Logger) Error(msg string, args ...interface{}) {
	logger.Log(msg, LevelError, args...)
}

// Helper function to log at debug level
func (logger *Logger) Fatal(msg string, args ...interface{}) {
	logger.Log(msg, LevelFatal, args...)
}
