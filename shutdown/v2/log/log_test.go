package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type logSink struct {
	logs []string
}

func (l *logSink) Write(p []byte) (n int, err error) {
	l.logs = append(l.logs, string(p))
	return len(p), nil
}

func (l *logSink) Index(i int) string {
	return l.logs[i]
}

func Test_shutdown_handler_does_log_when_service_will_be_stopped(t *testing.T) {
	logs := &logSink{}
	logger := zerolog.New(logs)
	handler := ShutdownLog{Logger: logger}

	handler.ServiceWillBeStopped("service name")

	assert.Len(t, logs.logs, 1)
	assert.Contains(t, logs.logs[0], "Stopping service name ...")
}

func Test_shutdown_handler_does_log_when_service_was_stopped_without_error(t *testing.T) {
	logs := &logSink{}
	logger := zerolog.New(logs)
	handler := ShutdownLog{Logger: logger}

	handler.ServiceWasStopped("service name")

	assert.Len(t, logs.logs, 1)
	assert.Contains(t, logs.logs[0], "service name stopped.")
}

func Test_shutdown_handler_does_log_when_service_was_stopped_with_error(t *testing.T) {
	logs := &logSink{}
	logger := zerolog.New(logs)
	handler := ShutdownLog{Logger: logger}

	handler.ServiceWasStopped("service name", fmt.Errorf("with error"))

	assert.Len(t, logs.logs, 1)
	assert.Contains(t, logs.logs[0], "Failed stopping 'service name'")
	assert.Contains(t, logs.logs[0], "\"no_alert\":true")
}
