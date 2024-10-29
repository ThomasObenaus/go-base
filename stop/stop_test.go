package stop

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
)

func Test_all_stoppable_items_are_stopped_in_order_given(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable1.EXPECT().String().AnyTimes()
	stoppable2.EXPECT().String().AnyTimes()
	stoppable3.EXPECT().String().AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop(),
		stoppable2.EXPECT().Stop(),
		stoppable1.EXPECT().Stop())

	// WHEN
	stoppableList.StopAllInOrder(zerolog.Logger{})
}

func Test_listener_is_called_when_a_service_is_about_to_be_stopped(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop(),
		stoppable2.EXPECT().Stop(),
		stoppable1.EXPECT().Stop(),
	)

	// WHEN
	stoppableList.StopAllInOrder(zerolog.Logger{})
}

func Test_listener_is_called_when_a_service_was_stopped(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop(),
		stoppable2.EXPECT().Stop(),
		stoppable1.EXPECT().Stop(),
	)

	// WHEN
	stoppableList.StopAllInOrder(zerolog.Logger{})
}

func Test_listener_is_called_when_a_service_could_not_be_stopped_without_error(t *testing.T) {
	// GIVEN
	stoppableList, ctrl, stoppable1, stoppable2, stoppable3 := createDefaultStopScenario2(t)
	defer ctrl.Finish()

	// IGNORE
	stoppable3.EXPECT().String().Return("service 3").AnyTimes()
	stoppable2.EXPECT().String().Return("service 2").AnyTimes()
	stoppable1.EXPECT().String().Return("service 1").AnyTimes()

	// EXPECT
	gomock.InOrder(
		stoppable3.EXPECT().Stop().Return(fmt.Errorf("error 3")),
		stoppable2.EXPECT().Stop().Return(fmt.Errorf("error 2")),
		stoppable1.EXPECT().Stop().Return(fmt.Errorf("error 1")),
	)

	// WHEN
	stoppableList.StopAllInOrder(zerolog.Logger{})
}

func createDefaultStopScenario2(t *testing.T) (Registry, *gomock.Controller, *MockStoppable, *MockStoppable, *MockStoppable) {
	mockCtrl := gomock.NewController(t)
	stoppable1 := NewMockStoppable(mockCtrl)
	stoppable2 := NewMockStoppable(mockCtrl)
	stoppable3 := NewMockStoppable(mockCtrl)

	return Registry{
		items: []Stoppable{stoppable3, stoppable2, stoppable1},
	}, mockCtrl, stoppable1, stoppable2, stoppable3
}
