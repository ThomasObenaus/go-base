package lifecycle

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func InstallShutdownHandler(orderedStopables []Stopable, logger zerolog.Logger) {
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)
	go shutdownHandler(shutDownChan, orderedStopables, logger)
	logger.Info().Msgf("Shutdown Handler installed for %d Stopables", len(orderedStopables))
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func shutdownHandler(shutdownChan <-chan os.Signal, orderedStopables []Stopable, logger zerolog.Logger) {
	s := <-shutdownChan
	logger.Info().Msgf("Received %v. Shutting down...", s)

	// Stop all components
	stop(orderedStopables, logger)
}

// stop calls Stop() on all Stopable in the list as they are ordered.
func stop(orderedStopables []Stopable, logger zerolog.Logger) {
	for _, stopable := range orderedStopables {
		name := stopable.Name()
		logger.Debug().Msgf("Stopping %s ...", name)
		err := stopable.Stop()
		if err != nil {
			logger.Error().Err(err).Msg("Failed stopping '%s'")
			continue
		}
		logger.Info().Msgf("%s stopped.", name)
	}
}
