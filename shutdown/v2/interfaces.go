package v2

import (
	stop2 "github.com/ThomasObenaus/go-base/shutdown/v2/stop"
)

// TODO: how to make a synchronized structure more visible
type stopIF interface {
	AddToFront(stoppable stop2.Stoppable) error
	AddToBack(stoppable1 stop2.Stoppable) error
	StopAllInOrder(listener stop2.Listener)
}

type signalHandlerIF interface {
	WaitForSignal()
	StopWaitingAndNotifyListener()
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
