package v2

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_can_register_a_stoppable_in_front(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	list := NewMockstopIF(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{stoppableItems: list}

	// EXPECT
	list.EXPECT().AddToFront(stoppable)

	// WHEN
	shutdownHandler.Register(stoppable, true)
}

func Test_logs_failure_if_stoppable_can_not_be_added_in_front(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStop := NewMockstopIF(mockCtrl)
	mockLog := NewMocklogIF(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockStop,
		log:            mockLog,
	}

	// EXPECT
	mockStop.EXPECT().AddToFront(stoppable).Return(fmt.Errorf("some error"))
	stoppable.EXPECT().String().Return("some service")
	mockLog.EXPECT().LogCanNotAddService("some service")

	// WHEN
	shutdownHandler.Register(stoppable, true)
}

func Test_can_register_a_stoppable_at_back(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	list := NewMockstopIF(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{stoppableItems: list}

	// EXPECT
	list.EXPECT().AddToBack(stoppable)

	// WHEN
	shutdownHandler.Register(stoppable)
}

func Test_logs_failure_if_stoppable_can_not_be_added_to_back(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStop := NewMockstopIF(mockCtrl)
	mockLog := NewMocklogIF(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockStop,
		log:            mockLog,
	}

	// EXPECT
	mockStop.EXPECT().AddToBack(stoppable).Return(fmt.Errorf("some error"))
	stoppable.EXPECT().String().Return("some service")
	mockLog.EXPECT().LogCanNotAddService("some service")

	// WHEN
	shutdownHandler.Register(stoppable)
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
	shutdownHandler.WaitForSignal()
}

func Test_can_stop(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSignalHandler := NewMocksignalHandlerIF(mockCtrl)

	shutdownHandler := ShutdownHandler{signalHandler: mockSignalHandler}

	// EXPECT
	mockSignalHandler.EXPECT().StopWaitingAndNotifyListener()

	// WHEN
	shutdownHandler.Stop()
}

func Test_logs_all_stop_related_events(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLog := NewMocklogIF(mockCtrl)
	mockStop := NewMockstopIF(mockCtrl)
	mockHealth := NewMockhealthIF(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockStop,
		log:            mockLog,
		health:         mockHealth,
	}

	// IGNORE
	mockHealth.EXPECT().ShutdownSignalReceived().AnyTimes()

	// EXPECT
	gomock.InOrder(
		mockLog.EXPECT().ShutdownSignalReceived(),
		mockStop.EXPECT().StopAllInOrder(mockLog),
	)

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
}

func Test_notifies_health_monitor_on_service_stop(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLog := NewMocklogIF(mockCtrl)
	mockStop := NewMockstopIF(mockCtrl)
	mockHealth := NewMockhealthIF(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockStop,
		log:            mockLog,
		health:         mockHealth,
	}

	// IGNORE
	mockLog.EXPECT().ShutdownSignalReceived().AnyTimes()
	mockStop.EXPECT().StopAllInOrder(gomock.Any()).AnyTimes()

	// EXPECT
	mockHealth.EXPECT().ShutdownSignalReceived()

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
}

func Test_uses_health_monitor_to_report_health_status(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHealth := NewMockhealthIF(mockCtrl)
	shutdownHandler := ShutdownHandler{
		health: mockHealth,
	}

	// EXPECT
	gomock.InOrder(
		mockHealth.EXPECT().IsHealthy().Return(nil),
		mockHealth.EXPECT().IsHealthy().Return(fmt.Errorf("some error")),
		mockHealth.EXPECT().String().Return("some status"),
	)

	// WHEN
	err := shutdownHandler.IsHealthy()
	assert.NoError(t, err)

	err = shutdownHandler.IsHealthy()
	assert.Error(t, err)

	status := shutdownHandler.String()
	assert.Equal(t, "some status", status)
}

func Test_informs_stop_that_it_should_stop(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLog := NewMocklogIF(mockCtrl)
	mockStop := NewMockstopIF(mockCtrl)
	mockHealth := NewMockhealthIF(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockStop,
		log:            mockLog,
		health:         mockHealth,
	}

	// IGNORE
	mockHealth.EXPECT().ShutdownSignalReceived().AnyTimes()
	mockLog.EXPECT().ShutdownSignalReceived()

	// EXPECT
	mockStop.EXPECT().StopAllInOrder(mockLog)

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
}
