package v2

import (
	health2 "github.com/ThomasObenaus/go-base/shutdown/health"
	log2 "github.com/ThomasObenaus/go-base/shutdown/log"
	"github.com/ThomasObenaus/go-base/shutdown/signal"
	stop2 "github.com/ThomasObenaus/go-base/shutdown/stop"
	"github.com/rs/zerolog"
)

type ShutdownHandler struct {
	stoppableItems stopIF
	signalHandler  signalHandlerIF
	log            logIF
	health         healthIF
}

// TODO: how to test this
func InstallHandler(orderedStopables []stop2.Stoppable, logger zerolog.Logger) *ShutdownHandler {
	shutdownHandler := &ShutdownHandler{
		stoppableItems: &stop2.OrderedStoppableList{},
		log:            log2.ShutdownLog{Logger: logger},
		health:         &health2.Health{},
	}

	for _, stopable := range orderedStopables {
		err := shutdownHandler.stoppableItems.AddToBack(stopable)
		if err != nil {
			logger.Error().Err(err).Msgf("can not stopre Stoppable")
			return nil
		}
	}

	handler := signal.NewDefaultSignalHandler(shutdownHandler)
	shutdownHandler.signalHandler = handler

	return shutdownHandler
}

func (h *ShutdownHandler) Register(stoppable stop2.Stoppable, front ...bool) {
	addToBack := isFirstBoolUndefinedOrFalse(front)

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

func isFirstBoolUndefinedOrFalse(front []bool) bool {
	addToBack := true

	if len(front) > 0 {
		if front[0] {
			addToBack = false
		}
	}
	return addToBack
}

func (h *ShutdownHandler) WaitUntilSignal() {
	h.signalHandler.WaitForSignal()
}

func (h *ShutdownHandler) Stop() {
	h.signalHandler.StopWaitingAndNotifyListener()
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
