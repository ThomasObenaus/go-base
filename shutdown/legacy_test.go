package shutdown

import (
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
)

func Test_can_register_a_stoppable_in_front(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	list := NewMockstopIF(mockCtrl)
	stoppable1 := NewMockStoppable(mockCtrl)
	stoppable2 := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{registry: list}

	// EXPECT
	list.EXPECT().AddToFront(stoppable1)
	list.EXPECT().AddToFront(stoppable2)

	// WHEN
	shutdownHandler.Register(stoppable1, true)
	shutdownHandler.Register(stoppable2)
}

func Test_can_register_a_stoppable_at_back(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	list := NewMockstopIF(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{registry: list}

	// EXPECT
	list.EXPECT().AddToBack(stoppable)

	// WHEN
	shutdownHandler.Register(stoppable, false)
}

func Test_can_wait_for_signal(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSignalHandler := NewMocksignalHandlerIF(mockCtrl)

	shutdownHandler := ShutdownHandler{signalHandler: mockSignalHandler}

	// EXPECT
	mockSignalHandler.EXPECT().WaitForSignal()

	// WHEN
	shutdownHandler.WaitUntilSignal()
}

func Test_informs_signal_handler_when_to_stop_waiting_for_signal_and_shutdown(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSignalHandler := NewMocksignalHandlerIF(mockCtrl)

	shutdownHandler := ShutdownHandler{signalHandler: mockSignalHandler}

	// EXPECT
	mockSignalHandler.EXPECT().NotifyListenerAndStopWaiting()

	// WHEN
	shutdownHandler.ShutdownAllAndStopWaiting()
}

func Test_informs_stop_that_it_should_stop(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.Nop()
	mockStop := NewMockstopIF(mockCtrl)
	shutdownHandler := ShutdownHandler{
		registry: mockStop,
		logger:   logger,
	}

	// EXPECT
	mockStop.EXPECT().StopAllInOrder(logger)

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
}
