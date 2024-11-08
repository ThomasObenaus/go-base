package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	go handler.waitForSignalAndCallListener(signalChannel, listener)

	return &handler
}

func (h *Handler) waitForSignalAndCallListener(signalChannel chan os.Signal, listener Listener) {
	defer h.wg.Done()
	h.wg.Add(1)
	_, _ = <-signalChannel
	listener.ShutdownSignalReceived()
}

func (h *Handler) WaitForSignal() {
	time.Sleep(time.Millisecond * 20)
	h.wg.Wait()
}

func (h *Handler) NotifyListenerAndStopWaiting() {
	close(h.signalChannel)
}
