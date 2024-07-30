package log

import (
	"github.com/rs/zerolog"
)

type ShutdownLog struct {
	Logger zerolog.Logger
}

func (h ShutdownLog) ShutdownSignalReceived() {
	h.Logger.Info().Msgf("Received %v. Shutting down...", h)
}

func (h ShutdownLog) ServiceWillBeStopped(name string) {
	h.Logger.Debug().Msgf("Stopping %s ...", name)
}

func (h ShutdownLog) ServiceWasStopped(name string, err ...error) {
	hasError := len(err) > 0 && err[0] != nil

	if hasError {
		h.Logger.Error().Err(err[0]).Bool("no_alert", true).Msgf("Failed stopping '%s'", name)
		return
	}
	h.Logger.Info().Msgf("%s stopped.", name)
}

func (h ShutdownLog) LogCanNotAddService(serviceName string) {
	h.Logger.Error().Msgf("can not add service '%s' to shutdown list while shutting down", serviceName)
}
