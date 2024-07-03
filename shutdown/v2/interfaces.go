package v2

import (
	"github.com/ThomasObenaus/go-base/shutdown/v2/stop"
)

// TODO: how to make a synchronized structure more visible
type SynchronizedList interface {
	AddToFront(stoppable stop.Stoppable)
	AddToBack(stoppable1 stop.Stoppable)
	GetItems() []stop.Stoppable
}

type SignalHandler interface {
	WaitForSignal()
	StopWaitingAndNotifyListener()
}

type Log interface {
	ShutdownSignalReceived()
	ServiceWillBeStopped(name string)
	ServiceWasStopped(name string, err ...error)
}

type Health interface {
	ShutdownSignalReceived()
	IsHealthy() error
	String() string
}
