package shutdown

import (
	"github.com/ThomasObenaus/go-base/shutdown/stop"
)

type stopIF interface {
	AddToFront(stoppable stop.Stoppable) error
	AddToBack(stoppable1 stop.Stoppable) error
	StopAllInOrder(listener stop.Listener)
}

type signalHandlerIF interface {
	WaitForSignal()
}

type logIF interface {
	ShutdownSignalReceived()
	ServiceWillBeStopped(name string)
	ServiceWasStopped(name string, err ...error)
	LogCanNotAddService(serviceName string)
}

type healthIF interface {
	ShutdownSignalReceived()
	IsHealthy() error
	String() string
}
