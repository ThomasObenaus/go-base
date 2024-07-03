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

	list := NewMocksynchronizedList(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{stoppableItems: list}

	// EXPECT
	list.EXPECT().AddToFront(stoppable)

	// WHEN
	shutdownHandler.Register(stoppable, true)
}

func Test_can_register_a_stoppable_at_back(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	list := NewMocksynchronizedList(mockCtrl)
	stoppable := NewMockStoppable(mockCtrl)
	shutdownHandler := ShutdownHandler{stoppableItems: list}

	// EXPECT
	list.EXPECT().AddToBack(stoppable)

	// WHEN
	shutdownHandler.Register(stoppable)
}

func Test_can_wait_for_signal(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSignalHandler := NewMocksignalHandler(mockCtrl)

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

	mockSignalHandler := NewMocksignalHandler(mockCtrl)

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

	mockLog := NewMocklog(mockCtrl)
	mockItemList := NewMocksynchronizedList(mockCtrl)
	mockHealth := NewMockhealth(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockItemList,
		log:            mockLog,
		health:         mockHealth,
	}

	// IGNORE
	mockItemList.EXPECT().GetItems().AnyTimes()
	mockHealth.EXPECT().ShutdownSignalReceived().AnyTimes()

	// EXPECT
	gomock.InOrder(
		mockLog.EXPECT().ShutdownSignalReceived(),
		mockLog.EXPECT().ServiceWillBeStopped("some service"),
		mockLog.EXPECT().ServiceWasStopped("some service", fmt.Errorf("some error")),
	)

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
	shutdownHandler.ServiceWillBeStopped("some service")
	shutdownHandler.ServiceWasStopped("some service", fmt.Errorf("some error"))
}

func Test_notifies_health_monitor_on_service_stop(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLog := NewMocklog(mockCtrl)
	mockItemList := NewMocksynchronizedList(mockCtrl)
	mockHealth := NewMockhealth(mockCtrl)
	shutdownHandler := ShutdownHandler{
		stoppableItems: mockItemList,
		log:            mockLog,
		health:         mockHealth,
	}

	// IGNORE
	mockItemList.EXPECT().GetItems().AnyTimes()
	mockLog.EXPECT().ShutdownSignalReceived().AnyTimes()

	// EXPECT
	mockHealth.EXPECT().ShutdownSignalReceived()

	// WHEN
	shutdownHandler.ShutdownSignalReceived()
}

func Test_uses_health_monitor_to_report_health_status(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHealth := NewMockhealth(mockCtrl)
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
