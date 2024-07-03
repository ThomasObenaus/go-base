package v2

import (
	"github.com/ThomasObenaus/go-base/shutdown"
	"github.com/ThomasObenaus/go-base/shutdown/v2/signal"
	"github.com/ThomasObenaus/go-base/shutdown/v2/stop"
	"github.com/ThomasObenaus/go-base/shutdown/v2/stop/list"
	"github.com/rs/zerolog"
)

type ShutdownHandler struct {
	stoppableItems synchronizedList
	signalHandler  signalHandler
	log            log
	health         health
}

func NewLegacyShutdownHandler(logger zerolog.Logger) *ShutdownHandler {
	shutdownHandler := &ShutdownHandler{
		stoppableItems: &list.SynchronizedList{},
		log:            log.ShutdownLog{Logger: logger},
		health:         &health.Health{},
	}

	handler := signal.NewDefaultSignalHandler(shutdownHandler)
	shutdownHandler.signalHandler = handler

	return shutdownHandler
}

func (h *ShutdownHandler) Register(stoppable shutdown.Stopable, front ...bool) {
	if len(front) > 0 {
		if front[0] {
			h.stoppableItems.AddToFront(stoppable)
			return
		}
	}

	h.stoppableItems.AddToBack(stoppable)
}

func (h *ShutdownHandler) WaitForSignal() {
	h.signalHandler.WaitForSignal()
}

func (h *ShutdownHandler) Stop() {
	h.signalHandler.StopWaitingAndNotifyListener()
}

func (h *ShutdownHandler) ShutdownSignalReceived() {
	h.log.ShutdownSignalReceived()
	h.health.ShutdownSignalReceived()
	stop.Stop(h.stoppableItems.GetItems(), h)
}

func (h *ShutdownHandler) ServiceWillBeStopped(name string) {
	h.log.ServiceWillBeStopped(name)
}

func (h *ShutdownHandler) ServiceWasStopped(name string, err ...error) {
	h.log.ServiceWasStopped(name, err...)
}

func (h *ShutdownHandler) IsHealthy() error {
	return h.health.IsHealthy()
}

func (h *ShutdownHandler) String() string {
	return h.health.String()
}
