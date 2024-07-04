package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Handler struct {
	signalChannel chan os.Signal
	wg            sync.WaitGroup
}

type Listener interface {
	ShutdownSignalReceived()
}

func NewDefaultSignalHandler(listener Listener) *Handler {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	return NewSignalHandler(signals, listener)
}

func NewSignalHandler(signalChannel chan os.Signal, listener Listener) *Handler {
	handler := Handler{signalChannel: signalChannel}

	go func() {
		handler.waitForSignalAndCallListener(signalChannel, listener)
	}()

	return &handler
}

func (h *Handler) waitForSignalAndCallListener(signalChannel chan os.Signal, listener Listener) {
	h.wg.Add(1)
	_, _ = <-signalChannel
	listener.ShutdownSignalReceived()
	h.wg.Done()
}

func (h *Handler) WaitForSignal() {
	h.wg.Wait()
}

func (h *Handler) NotifyListenerAndStopWaiting() {
	close(h.signalChannel)
}
