package shutdown

import (
	"github.com/ThomasObenaus/go-base/shutdown/stop"
	"github.com/rs/zerolog"
)

type stopIF interface {
	AddToFront(stoppable stop.Stoppable) error
	AddToBack(stoppable1 stop.Stoppable) error
	StopAllInOrder(logger zerolog.Logger)
}

type signalHandlerIF interface {
	WaitForSignal()
	NotifyListenerAndStopWaiting()
}

type healthIF interface {
	ShutdownSignalReceived()
	IsHealthy() error
	String() string
}
