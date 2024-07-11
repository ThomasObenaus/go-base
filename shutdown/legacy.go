package shutdown

import (
	"github.com/ThomasObenaus/go-base/shutdown/health"
	"github.com/ThomasObenaus/go-base/shutdown/log"
	"github.com/ThomasObenaus/go-base/shutdown/signal"
	"github.com/ThomasObenaus/go-base/shutdown/stop"
	"github.com/rs/zerolog"
)

type ShutdownHandler struct {
	stoppableItems stopIF
	signalHandler  signalHandlerIF
	log            logIF
	health         healthIF
}

func InstallHandler(orderedStopables []stop.Stoppable, logger zerolog.Logger) *ShutdownHandler {
	shutdownHandler := &ShutdownHandler{
		stoppableItems: &stop.OrderedStoppableList{},
		log:            log.ShutdownLog{Logger: logger},
		health:         &health.Health{},
	}

	for _, stopable := range orderedStopables {
		err := shutdownHandler.stoppableItems.AddToBack(stopable)
		if err != nil {
			logger.Error().Err(err).Msgf("unexpected error adding stoppable to internal list")
			return nil
		}
	}

	handler := signal.NewDefaultSignalHandler(shutdownHandler)
	shutdownHandler.signalHandler = handler

	return shutdownHandler
}

func (h *ShutdownHandler) Register(stoppable stop.Stoppable, front ...bool) {
	addToBack := isEmptyOrFirstEntryFalse(front)

	if addToBack {
		err := h.stoppableItems.AddToBack(stoppable)
		if err != nil {
			serviceName := stoppable.String()
			h.log.LogCanNotAddService(serviceName)
		}
		return
	}
	err := h.stoppableItems.AddToFront(stoppable)
	if err != nil {
		serviceName := stoppable.String()
		h.log.LogCanNotAddService(serviceName)
	}
}

func isEmptyOrFirstEntryFalse(list []bool) bool {
	if len(list) == 0 {
		return true
	}

	return !list[0]
}

func (h *ShutdownHandler) WaitUntilSignal() {
	h.signalHandler.WaitForSignal()
}

func (h *ShutdownHandler) ShutdownSignalReceived() {
	h.log.ShutdownSignalReceived()
	h.health.ShutdownSignalReceived()
	h.stoppableItems.StopAllInOrder(h.log)
}

func (h *ShutdownHandler) IsHealthy() error {
	return h.health.IsHealthy()
}

func (h *ShutdownHandler) String() string {
	return h.health.String()
}