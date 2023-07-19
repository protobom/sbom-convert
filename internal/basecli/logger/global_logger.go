package logger

import "github.com/sirupsen/logrus"

var log Logger

// InitializeGlobalLogger this function initializes the global logger with the logrus one based on the configuration.
func SetDefaultGlobalLogger(cfg Config) error {
	log, err := NewDefaultLogger(cfg)
	if err != nil {
		return err
	}
	SetGlobalLogger(log)
	return nil
}

// SetGlobalLogger sets the global logger
func SetGlobalLogger(logger Logger) {
	log = logger
}

// GetGlobalLogger gets the global logger
func GetGlobalLogger() Logger {
	return log
}

// Errorf formats and logs the data at error level
func Errorf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Errorf(format, args...)
}

// Error logs the data at error level
func Error(args ...interface{}) {
	if log == nil {
		return
	}
	log.Error(args...)
}

// Warnf formats the string and logs the data at warn level
func Warnf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Warnf(format, args...)
}

// Warn formats the string and logs the data at warn level
func Warn(args ...interface{}) {
	if log == nil {
		return
	}
	log.Warn(args...)
}

// Infof formats the data and logs the data at info level
func Infof(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Infof(format, args...)
}

// Info formats the data and logs the data at info level
func Info(args ...interface{}) {
	if log == nil {
		return
	}
	log.Info(args...)
}

// Debugf formats the data and logs the data at debug level
func Debugf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debugf(format, args...)
}

// Debug logs the data at debug level
func Debug(args ...interface{}) {
	if log == nil {
		return
	}
	log.Debug(args...)
}

// Debugf formats the data and logs the data at debug level
func Tracef(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Tracef(format, args...)
}

// Debug logs the data at debug level
func Trace(args ...interface{}) {
	if log == nil {
		return
	}
	log.Trace(args...)
}

// Debugf formats the data and logs the data at debug level
func Panicf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Panicf(format, args...)
}

// Debug logs the data at debug level
func Panic(args ...interface{}) {
	if log == nil {
		return
	}
	log.Panic(args...)
}

// SetLevel sets the logging level
func SetLevel(level logrus.Level) {
	if log == nil {
		return
	}
	log.SetLevel(level)
}

// GetLevel gets the loging level
func GetLevel() logrus.Level {
	return log.GetLevel()
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
