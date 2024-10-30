package stop

import (
	"errors"
	"github.com/rs/zerolog"
	"sync"
)

type Registry struct {
	items                        []Stoppable
	mux                          sync.Mutex
	shutdownInProgressOrComplete bool
}

func (l *Registry) AddToFront(stoppable Stoppable) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.shutdownInProgressOrComplete {
		return errors.New("can not add services while shutting down in progress")
	}

	l.items = append([]Stoppable{stoppable}, l.items...)
	return nil
}

func (l *Registry) AddToBack(stoppable1 Stoppable) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.shutdownInProgressOrComplete {
		return errors.New("can not add services while shutting down in progress")
	}

	l.items = append(l.items, stoppable1)

	return nil
}

func (l *Registry) StopAllInOrder(logger zerolog.Logger) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.shutdownInProgressOrComplete = true
	stop(l.items, logger)
}

func stop(stoppableItems []Stoppable, logger zerolog.Logger) {
	for _, stoppable := range stoppableItems {
		serviceName := stoppable.String()
		logger.Debug().Msgf("Stopping %s ...", serviceName)
		err := stoppable.Stop()
		if err != nil {
			logger.Error().Err(err).Bool("no_alert", true).Msgf("Failed stopping '%s'", serviceName)
			continue
		}
		logger.Info().Msgf("%s stopped.", serviceName)
	}
}
