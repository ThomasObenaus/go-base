package shutdown

import (
	"github.com/ThomasObenaus/go-base/signal"
	"github.com/ThomasObenaus/go-base/stop"
	"github.com/rs/zerolog"
	"sync/atomic"
)

type ShutdownHandler struct {
	logger            zerolog.Logger
	isShutdownPending atomic.Bool
	registry          stopIF
	signalHandler     signalHandlerIF
}

// InstallHandler installs a handler for syscall.SIGINT, syscall.SIGTERM
func InstallHandler(orderedStopables []stop.Stoppable, logger zerolog.Logger) *ShutdownHandler {
	shutdownHandler := &ShutdownHandler{
		registry: &stop.Registry{},
	}

	for _, stopable := range orderedStopables {
		err := shutdownHandler.registry.AddToBack(stopable)
		if err != nil {
			logger.Error().Err(err).Msgf("unexpected error adding stoppable to internal list")
			return nil
		}
	}

	handler := signal.NewDefaultSignalHandler(shutdownHandler)
	shutdownHandler.signalHandler = handler

	return shutdownHandler
}

// Register a Stopable for shutdown handling. Per default the Stopable
// is added to the front of the list of Stopable's this means the
// Stopable that was the last one registered will be the first being called for shutdown.
// If you call Register(stopable,false) you can add this Stopable to the end
// of the list of registered Stopables.
func (h *ShutdownHandler) Register(stoppable stop.Stoppable, front ...bool) {
	addToFront := isEmptyOrFirstEntryTrue(front)

	if addToFront {
		err := h.registry.AddToFront(stoppable)
		if err != nil {
			serviceName := stoppable.String()
			h.logger.Error().Msgf("can not add service '%s' to shutdown list while shutting down", serviceName)
		}
		return
	}
	err := h.registry.AddToBack(stoppable)
	if err != nil {
		serviceName := stoppable.String()
		h.logger.Error().Msgf("can not add service '%s' to shutdown list while shutting down", serviceName)
	}
}

func isEmptyOrFirstEntryTrue(list []bool) bool {
	if len(list) == 0 {
		return true
	}

	return list[0]
}

func (h *ShutdownHandler) WaitUntilSignal() {
	h.signalHandler.WaitForSignal()
}

func (h *ShutdownHandler) ShutdownAllAndStopWaiting() {
	h.signalHandler.NotifyListenerAndStopWaiting()
}

func (h *ShutdownHandler) ShutdownSignalReceived() {
	h.logger.Info().Msgf("Received %v. Shutting down...", h)
	h.isShutdownPending.Store(true)
	h.registry.StopAllInOrder(h.logger)
}
