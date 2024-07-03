package list

import (
	"github.com/ThomasObenaus/go-base/shutdown/v2/stop"
	"sync"
)

type SynchronizedList struct {
	items []stop.Stoppable
	mux   sync.Mutex
}

func (l *SynchronizedList) AddToFront(stoppable stop.Stoppable) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.items = append([]stop.Stoppable{stoppable}, l.items...)
}

func (l *SynchronizedList) AddToBack(stoppable1 stop.Stoppable) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.items = append(l.items, stoppable1)
}

func (l *SynchronizedList) GetItems() []stop.Stoppable {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.items
}
