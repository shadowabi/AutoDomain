package log

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const (
	LogLevelTrace = "trace"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"
)

// Init func is a function to init logrus with specific log level
func Init(level string) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(logFormat())
	log.SetLevel(logLevel(level))
}

// logLevel search level strings return correct Level
func logLevel(level string) log.Level {
	switch level {
	case LogLevelTrace:
		return log.TraceLevel
	case LogLevelDebug:
		return log.DebugLevel
	case LogLevelInfo:
		return log.InfoLevel
	case LogLevelWarn:
		return log.WarnLevel
	case LogLevelError:
		return log.ErrorLevel
	case LogLevelFatal:
		return log.FatalLevel
	case LogLevelPanic:
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

// logFormat sets log format by using prefixed "x-cray/logrus-prefixed-formatter"
func logFormat() log.Formatter {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})
	return formatter
}
