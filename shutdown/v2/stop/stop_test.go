package stop

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_all_stoppable_items_are_stopped_in_order_given(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, listener, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable1.EXPECT().String().AnyTimes()
	stoppable2.EXPECT().String().AnyTimes()
	stoppable3.EXPECT().String().AnyTimes()
	listener.EXPECT().ServiceWillBeStopped(gomock.Any()).AnyTimes()
	listener.EXPECT().ServiceWasStopped(gomock.Any()).AnyTimes()
	listener.EXPECT().ServiceWasStopped(gomock.Any(), gomock.Any()).AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop(),
		stoppable2.EXPECT().Stop(),
		stoppable1.EXPECT().Stop())

	// WHEN
	stoppableList.StopAllInOrder(listener)
}

func Test_listener_is_called_when_a_service_is_about_to_be_stopped(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, listener, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()
	listener.EXPECT().ServiceWasStopped(gomock.Any()).AnyTimes()
	listener.EXPECT().ServiceWasStopped(gomock.Any(), gomock.Any()).AnyTimes()

	// EXPECT
	gomock.InOrder(
		listener.EXPECT().ServiceWillBeStopped("service 3"),
		stoppable3.EXPECT().Stop(),
		listener.EXPECT().ServiceWillBeStopped("service 2"),
		stoppable2.EXPECT().Stop(),
		listener.EXPECT().ServiceWillBeStopped("service 1"),
		stoppable1.EXPECT().Stop(),
	)

	// WHEN
	stoppableList.StopAllInOrder(listener)
}

func Test_listener_is_called_when_a_service_was_stopped(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, listener, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()
	listener.EXPECT().ServiceWillBeStopped(gomock.Any()).AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop(),
		listener.EXPECT().ServiceWasStopped("service 3", gomock.Any()),
		stoppable2.EXPECT().Stop(),
		listener.EXPECT().ServiceWasStopped("service 2", gomock.Any()),
		stoppable1.EXPECT().Stop(),
		listener.EXPECT().ServiceWasStopped("service 1", gomock.Any()),
	)

	// WHEN
	stoppableList.StopAllInOrder(listener)
}

func Test_listener_is_called_when_a_service_could_not_be_stopped_without_error(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, listener, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()
	listener.EXPECT().ServiceWillBeStopped(gomock.Any()).AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop().Return(fmt.Errorf("error 3")),
		listener.EXPECT().ServiceWasStopped("service 3", fmt.Errorf("error 3")),
		stoppable2.EXPECT().Stop().Return(fmt.Errorf("error 2")),
		listener.EXPECT().ServiceWasStopped("service 2", fmt.Errorf("error 2")),
		stoppable1.EXPECT().Stop().Return(fmt.Errorf("error 1")),
		listener.EXPECT().ServiceWasStopped("service 1", fmt.Errorf("error 1")),
	)

	// WHEN
	stoppableList.StopAllInOrder(listener)
}

func createDefaultStopScenario2(t *testing.T) (OrderedStoppableList, *gomock.Controller, *MockListener, *MockStoppable, *MockStoppable, *MockStoppable) {
	mockCtrl := gomock.NewController(t)
	stoppable1 := NewMockStoppable(mockCtrl)
	stoppable2 := NewMockStoppable(mockCtrl)
	stoppable3 := NewMockStoppable(mockCtrl)

	listener := NewMockListener(mockCtrl)

	return OrderedStoppableList{
		items: []Stoppable{stoppable3, stoppable2, stoppable1},
	}, mockCtrl, listener, stoppable1, stoppable2, stoppable3
}
