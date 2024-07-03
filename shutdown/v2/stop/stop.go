package stop

func Stop(stoppableItems []Stoppable, listener Listener) {
	for _, stoppable := range stoppableItems {
		serviceName := stoppable.String()
		listener.ServiceWillBeStopped(serviceName)
		err := stoppable.Stop()
		listener.ServiceWasStopped(serviceName, err)
	}
}
