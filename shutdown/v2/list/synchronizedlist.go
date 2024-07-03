package list

import "sync"

type SynchronizedList[T interface{}] struct {
	items []T
	mux   sync.Mutex
}

func (l *SynchronizedList[T]) AddToFront(stoppable T) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.items = append([]T{stoppable}, l.items...)
}

func (l *SynchronizedList[T]) AddToBack(stoppable1 T) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.items = append(l.items, stoppable1)
}

func (l *SynchronizedList[T]) GetItems() []T {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.items
}
