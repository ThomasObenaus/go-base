package stop

import (
	"errors"
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

func (l *OrderedStoppableList) StopAllInOrder(listener Listener) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.isShuttingDown = true
	stop(l.items, listener)
}

func stop(stoppableItems []Stoppable, listener Listener) {
	for _, stoppable := range stoppableItems {
		serviceName := stoppable.String()
		listener.ServiceWillBeStopped(serviceName)
		err := stoppable.Stop()
		listener.ServiceWasStopped(serviceName, err)
	}
}
