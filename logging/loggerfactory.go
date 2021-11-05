package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// LoggerFactory is a factory that can be used to create named loggers using the same aligned configuration and namespace.
type LoggerFactory interface {
	NewNamedLogger(name string) zerolog.Logger
	Level() zerolog.Level
	IsStructuredLogging() bool
}

// Option is the struct for defining optional parameters for LoggerFactory
type Option func(*loggerFactoryImpl)

// Level sets the log-level used for the loggers created through this factory
func Level(level zerolog.Level) Option {
	return func(lf *loggerFactoryImpl) {
		lf.logLevel = level
	}
}

// New creates a new LoggerFactory which then can be used to create configured named loggers (log channels)
func New(structuredLogging, unixTimeStamp, disableColoredLogs bool, options ...Option) LoggerFactory {

	factory := &loggerFactoryImpl{
		structuredLogging:  structuredLogging,
		disableColoredLogs: disableColoredLogs,
		logLevel:           zerolog.DebugLevel,
	}

	for _, opt := range options {
		opt(factory)
	}

	// default format for the timestamp
	factory.timeFieldFormat = time.StampMilli //time.RFC3339

	if unixTimeStamp {
		// UNIX Time is faster and smaller than most timestamps
		// If you set zerolog.TimeFieldFormat to an empty string,
		// logs will write with UNIX time
		factory.timeFieldFormat = zerolog.TimeFormatUnix
	}

	return factory
}

type loggerFactoryImpl struct {
	structuredLogging  bool
	disableColoredLogs bool
	logLevel           zerolog.Level
	timeFieldFormat    string // for allowed values please refer to zerolog.TimeFieldFormat
}

// NewNamedLogger creates a new named and configured log-channel (logger)
func (lf *loggerFactoryImpl) NewNamedLogger(name string) zerolog.Logger {

	if lf.structuredLogging {
		return zerolog.New(os.Stdout).Level(lf.logLevel).With().Timestamp().Str("logger", name).Logger()
	}

	return zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		NoColor: lf.disableColoredLogs, Out: os.Stderr,
		TimeFormat: lf.timeFieldFormat,
	}).Level(lf.logLevel).With().Timestamp().Str("logger", name).Logger()
}

func (lf *loggerFactoryImpl) IsStructuredLogging() bool {
	return lf.structuredLogging
}

func (lf *loggerFactoryImpl) Level() zerolog.Level {
	return lf.logLevel
}
