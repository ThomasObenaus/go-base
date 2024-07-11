package signal

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_can_create_signal_handler_which_responds_to_actual_signals(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	listener := NewMockListener(mockCtrl)
	done := make(chan struct{})

	handler := NewDefaultSignalHandler(listener)
	assert.NotNil(t, handler)

	// EXPECT
	listener.EXPECT().ShutdownSignalReceived().Do(func() {
		close(done)
	})

	// WHEN
	go func() {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.NoError(t, err)
	}()

	// THEN

	timeout := time.After(time.Second)
	select {
	case <-done:
	case <-timeout:
		t.Errorf("signal handler listener was never called")
	}
}

func Test_can_create_signal_handler_which_calls_listener_when_signal_is_received(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	listener := NewMockListener(mockCtrl)
	done := make(chan struct{})

	shutDownChan := make(chan os.Signal, 1)

	handler := NewSignalHandler(shutDownChan, listener)
	assert.NotNil(t, handler)

	// EXPECT
	listener.EXPECT().ShutdownSignalReceived().Do(func() {
		close(done)
	})

	// WHEN
	go func() {
		shutDownChan <- syscall.SIGTERM
	}()

	// THEN

	timeout := time.After(time.Second)
	select {
	case <-done:
	case <-timeout:
		t.Errorf("signal handler listener was never called")
	}
}

func Test_signal_handler_will_not_block_if_signal_is_received(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	listener := NewMockListener(mockCtrl)

	done := make(chan struct{})

	shutDownChan := make(chan os.Signal, 1)
	signalHandler := NewSignalHandler(shutDownChan, listener)

	// IGNORE
	listener.EXPECT().ShutdownSignalReceived().AnyTimes()

	// WHEN
	go func() {
		shutDownChan <- syscall.SIGTERM
	}()

	// THEN
	go func() {
		signalHandler.WaitForSignal()
		close(done)
	}()

	timeout := time.After(time.Second)
	select {
	case <-done:
	case <-timeout:
		t.Errorf("wait for signal blocked longer then expected")
	}
}

func Test_signal_handler_will_block_if_no_signal_is_received(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	listener := NewMockListener(mockCtrl)

	done := make(chan struct{})

	shutDownChan := make(chan os.Signal, 1)
	signalHandler := NewSignalHandler(shutDownChan, listener)

	// IGNORE
	listener.EXPECT().ShutdownSignalReceived().AnyTimes()

	// WHEN - no signal is send

	// THEN
	go func() {
		signalHandler.WaitForSignal()
		close(done)
	}()

	timeout := time.After(time.Millisecond * 20)
	select {
	case <-done:
		t.Errorf("wait for signal did not block")
	case <-timeout:
	}
}
