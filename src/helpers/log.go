package helpers

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logging

type logging struct {
	errorLogger   *logrus.Logger
	successLogger *logrus.Logger
}

func (l *logging) Panic(msg string) {
	l.errorLogger.Panic(msg)
}

func (l *logging) Fatal(msg string) {
	l.errorLogger.Fatal(msg)
}

func (l *logging) Error(msg string) {
	l.errorLogger.Error(msg)
}

func (l *logging) Warn(msg string) {
	l.successLogger.Warn(msg)
}

func (l *logging) Info(msg string) {
	l.successLogger.Info(msg)
}

func (l *logging) Debug(msg string) {
	l.successLogger.Debug(msg)
}

func (l *logging) Trace(msg string) {
	l.successLogger.Trace(msg)
}

func getLogLevel(level uint32) logrus.Level {
	switch level {
	case 0:
		return logrus.PanicLevel
	case 1:
		return logrus.FatalLevel
	case 2:
		return logrus.ErrorLevel
	case 3:
		return logrus.WarnLevel
	case 4:
		return logrus.InfoLevel
	case 5:
		return logrus.DebugLevel
	default:
		return logrus.TraceLevel
	}
}
func NewLogger(
	errorLogPath string,
	successLogPath string,
	logLevel uint32,
	rotationLogTime int,
	rotationLogSize int,
	rotationLogBackups int,
) *logging {
	if Logger != nil {
		return Logger
	}
	// log files
	errorLogFile := &lumberjack.Logger{
		Filename:   errorLogPath,
		MaxSize:    rotationLogSize, // megabytes
		MaxBackups: rotationLogBackups,
		Compress:   true,
	}
	successLogFile := &lumberjack.Logger{
		Filename:   successLogPath,
		MaxSize:    rotationLogSize, // megabytes
		MaxBackups: rotationLogBackups,
		Compress:   true,
	}

	// loggers
	errorLogger := logrus.New()
	errorLogger.SetOutput(errorLogFile)
	errorLogger.SetLevel(logrus.ErrorLevel)

	successLogger := logrus.New()
	successLogger.SetOutput(successLogFile)
	successLogger.SetLevel(getLogLevel(logLevel))

	Logger = &logging{errorLogger: errorLogger, successLogger: successLogger}
	return Logger
}
