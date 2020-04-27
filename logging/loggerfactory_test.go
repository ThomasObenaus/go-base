package logging

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	lf := New(false, false, false)
	assert.NotNil(t, lf)
	assert.NotEmpty(t, zerolog.TimeFieldFormat)

	lf = New(false, true, false)
	assert.NotNil(t, lf)
	assert.Empty(t, zerolog.TimeFieldFormat)
}

func TestNewNamedLogger(t *testing.T) {

	loggerFactory := New(true, false, false)

	logger := loggerFactory.NewNamedLogger("MyTestLogger")
	strout := strings.Builder{}
	loggerDup := logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger")
	strout.Reset()

	loggerFactory = New(false, false, false)

	logger = loggerFactory.NewNamedLogger("MyTestLogger2")
	loggerDup = logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger2")
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
}

func TestLoglevel(t *testing.T) {

	loggerFactory := New(true, false, false, Level(zerolog.ErrorLevel))

	logger := loggerFactory.NewNamedLogger("MyTestLogger")
	strout := strings.Builder{}
	loggerDup := logger.Output(&strout)
	loggerDup.Info().Msg("INFO-MESSAGE")
	loggerDup.Error().Msg("ERROR-MESSAGE")
	assert.NotContains(t, strout.String(), "INFO-MESSAGE")
	assert.Contains(t, strout.String(), "ERROR-MESSAGE")
	assert.Equal(t, zerolog.ErrorLevel, logger.GetLevel())
	strout.Reset()

}

func ExampleNew() {
	// create the factory
	loggingFactory := New(true, false, false)

	// create new named logger
	logger := loggingFactory.NewNamedLogger("MyLogger")
	logger.Info().Msg("Hello World")

}
