package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/mishankov/go-utlz/cliutils"
)

type Logger struct {
	name string
	// logLevel LogLevel
	parent *Logger
}

func NewLogger(name string) Logger {
	return Logger{name: name}
}

func NewLoggerFromParent(name string, parent *Logger) Logger {
	return Logger{name: name, parent: parent}
}

func (l *Logger) FullLoggerName() string {
	if l.parent == nil {
		return l.name
	}

	return l.parent.FullLoggerName() + "." + l.name
}

func (l *Logger) GetLogLevelFromEnv() LogLevel {
	currentLoggerLevel := cliutils.GetEnvOrDefault("LOG_LEVEL_"+l.FullLoggerName(), "None")

	switch currentLoggerLevel {
	case "None":
	case logLevels.Debug.name:
		return logLevels.Debug
	case logLevels.Info.name:
		return logLevels.Info
	case logLevels.Warn.name:
		return logLevels.Warn
	case logLevels.Error.name:
		return logLevels.Error
	case logLevels.Fatal.name:
		return logLevels.Fatal
	}

	globalLoggerLevel := cliutils.GetEnvOrDefault("LOG_LEVEL", "None")

	switch globalLoggerLevel {
	case "None":
	case logLevels.Debug.name:
		return logLevels.Debug
	case logLevels.Info.name:
		return logLevels.Info
	case logLevels.Warn.name:
		return logLevels.Warn
	case logLevels.Error.name:
		return logLevels.Error
	case logLevels.Fatal.name:
		return logLevels.Fatal
	}

	return logLevels.Info
}

func (l *Logger) ShouldWriteLog(logLevel LogLevel) bool {
	return logLevel.level >= l.GetLogLevelFromEnv().level
}

type LogLevel struct {
	name  string
	level int
}

type LogLevels struct {
	Debug LogLevel
	Info  LogLevel
	Warn  LogLevel
	Error LogLevel
	Fatal LogLevel
}

var logLevels LogLevels = LogLevels{
	Debug: LogLevel{"Debug", 0},
	Info:  LogLevel{"Info", 1},
	Warn:  LogLevel{"Warn", 2},
	Error: LogLevel{"Error", 3},
	Fatal: LogLevel{"Fatal", 4},
}

func (l *Logger) Logf(logLevel LogLevel, message string, a ...any) {
	message = fmt.Sprintf(message, a...)
	l.Log(logLevel, message)
}

func (l *Logger) Log(logLevel LogLevel, message any) {
	if l.ShouldWriteLog(logLevel) {
		fmt.Printf("[%v] [%v] [%v] - %v\n", time.Now().Format("2006-01-02 15:04:05 GMT-0700"), l.FullLoggerName(), logLevel.name, message)
	}
}

func (l *Logger) Debug(message any) {
	l.Log(logLevels.Debug, message)
}
func (l *Logger) Debugf(message string, a ...any) {
	l.Logf(logLevels.Debug, message, a...)
}

func (l *Logger) Info(message any) {
	l.Log(logLevels.Info, message)
}
func (l *Logger) Infof(message string, a ...any) {
	l.Logf(logLevels.Info, message, a...)
}

func (l *Logger) Warn(message any) {
	l.Log(logLevels.Warn, message)
}
func (l *Logger) Warnf(message string, a ...any) {
	l.Logf(logLevels.Warn, message, a...)
}

func (l *Logger) Error(message any) {
	l.Log(logLevels.Error, message)
}
func (l *Logger) Errorf(message string, a ...any) {
	l.Logf(logLevels.Error, message, a...)
}

func (l *Logger) Fatal(message any) {
	l.Log(logLevels.Fatal, message)
	os.Exit(1)
}
func (l *Logger) Fatalf(message string, a ...any) {
	l.Logf(logLevels.Fatal, message, a...)
	os.Exit(1)
}
