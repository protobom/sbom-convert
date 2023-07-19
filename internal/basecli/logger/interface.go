package logger

import "github.com/sirupsen/logrus"

// Logger interface that helps logging
type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	SetLevel(level logrus.Level)
	GetLevel() logrus.Level
	WithFields(fields logrus.Fields) *logrus.Entry
}
