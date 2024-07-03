package stop

import "github.com/ThomasObenaus/go-base/shutdown"

type Listener interface {
	ServiceWillBeStopped(name string)
	ServiceWasStopped(name string, err ...error)
}

func Stop(stoppableItems []shutdown.Stopable, listener Listener) {
	for _, stoppable := range stoppableItems {
		serviceName := stoppable.String()
		listener.ServiceWillBeStopped(serviceName)
		err := stoppable.Stop()
		listener.ServiceWasStopped(serviceName, err)
	}
}
