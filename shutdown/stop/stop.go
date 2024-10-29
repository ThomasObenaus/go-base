package stop

import (
	"errors"
	"github.com/rs/zerolog"
	"sync"
)

type OrderedStoppableList struct {
	items          []Stoppable
	mux            sync.Mutex
	isShuttingDown bool
}

func (l *OrderedStoppableList) AddToFront(stoppable Stoppable) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.isShuttingDown {
		return errors.New("can not add services while shutting down in progress")
	}

	l.items = append([]Stoppable{stoppable}, l.items...)
	return nil
}

func (l *OrderedStoppableList) AddToBack(stoppable1 Stoppable) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.isShuttingDown {
		return errors.New("can not add services while shutting down in progress")
	}

	l.items = append(l.items, stoppable1)

	return nil
}

func (l *OrderedStoppableList) StopAllInOrder(logger zerolog.Logger) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.isShuttingDown = true
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
