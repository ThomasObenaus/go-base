package shutdown

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

// Handler represents a handler for shutdown events
type Handler struct {
	logger            zerolog.Logger
	isShutdownPending bool
}

// InstallHandler installs a handler for syscall.SIGINT, syscall.SIGTERM
func InstallHandler(orderedStopables []Stopable, logger zerolog.Logger) *Handler {
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)

	handler := &Handler{
		logger:            logger,
		isShutdownPending: false,
	}

	go handler.shutdownHandler(shutDownChan, orderedStopables, logger)
	handler.logger.Info().Msgf("Shutdown Handler installed for %d Stopables", len(orderedStopables))

	return handler
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func (h *Handler) shutdownHandler(shutdownChan <-chan os.Signal, orderedStopables []Stopable, logger zerolog.Logger) {
	s := <-shutdownChan
	h.isShutdownPending = true
	logger.Info().Msgf("Received %v. Shutting down...", s)

	// Stop all components
	stop(orderedStopables, logger)
}

// stop calls Stop() on all Stopable in the list as they are ordered.
func stop(orderedStopables []Stopable, logger zerolog.Logger) {
	for _, stopable := range orderedStopables {
		name := stopable.String()
		logger.Debug().Msgf("Stopping %s ...", name)
		err := stopable.Stop()
		if err != nil {
			logger.Error().Err(err).Msg("Failed stopping '%s'")
			continue
		}
		logger.Info().Msgf("%s stopped.", name)
	}
}
