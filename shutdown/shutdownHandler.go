package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

// Handler represents a handler for shutdown events
type Handler struct {
	logger            zerolog.Logger
	isShutdownPending bool
	wg                sync.WaitGroup
	orderedStopables  []Stopable
}

// InstallHandler installs a handler for syscall.SIGINT, syscall.SIGTERM
func InstallHandler(orderedStopables []Stopable, logger zerolog.Logger) *Handler {
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)

	handler := &Handler{
		logger:            logger,
		isShutdownPending: false,
		orderedStopables:  make([]Stopable, 0),
	}
	handler.orderedStopables = append(handler.orderedStopables, orderedStopables...)

	go handler.shutdownHandler(shutDownChan, logger)
	handler.logger.Info().Msgf("Shutdown Handler installed")
	return handler
}

// Register a Stopable for shutdown handling. Per default the Stopable
// is added to the front of the list of Stopable's this means the
// Stopable that was the last one registered will be the first being called for shutdown.
// If you call Register(stopable,false) you can add this Stopable to the end
// of the list of registered Stopables.
func (h *Handler) Register(stopable Stopable, front ...bool) {
	pushFront := true

	if len(front) > 0 {
		pushFront = front[0]
	}

	if pushFront {
		h.orderedStopables = append([]Stopable{stopable}, h.orderedStopables...)
	} else {
		h.orderedStopables = append(h.orderedStopables, stopable)
	}
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func (h *Handler) shutdownHandler(shutdownChan <-chan os.Signal, logger zerolog.Logger) {
	h.wg.Add(1)
	defer h.wg.Done()

	s := <-shutdownChan
	h.isShutdownPending = true
	logger.Info().Msgf("Received %v. Shutting down...", s)

	// Stop all components
	stop(h.orderedStopables, logger)
}

// WaitUntilSignal waits/ blocks until either syscall.SIGINT or syscall.SIGTERM was issued to the process
func (h *Handler) WaitUntilSignal() {
	time.Sleep(time.Millisecond * 20)
	h.wg.Wait()
}

// stop calls Stop() on all Stopable in the list as they are ordered.
func stop(orderedStopables []Stopable, logger zerolog.Logger) {
	for _, stopable := range orderedStopables {
		name := stopable.String()
		logger.Debug().Msgf("Stopping %s ...", name)
		err := stopable.Stop()
		if err != nil {
			logger.Error().Err(err).Msgf("Failed stopping '%s'", name)
			continue
		}
		logger.Info().Msgf("%s stopped.", name)
	}
}
