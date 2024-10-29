package shutdown

import (
	"fmt"
	"github.com/ThomasObenaus/go-base/shutdown/stop"
	"github.com/ThomasObenaus/go-base/signal"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type signalMock struct {
}

func (s signalMock) Signal() {
}
func (s signalMock) String() string {
	return ""
}

// These tests are the old ones from before the refactoring
// They confirm, that everything that was tested before still behaves the same

func Test_ShutdownHandler(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	items := &stop.OrderedStoppableList{}

	stopable1 := NewMockStoppable(mockCtrl)
	err := items.AddToFront(stopable1)
	require.NoError(t, err)
	stopable2 := NewMockStoppable(mockCtrl)
	err = items.AddToBack(stopable2)
	require.NoError(t, err)

	h := ShutdownHandler{
		logger:         zerolog.Nop(),
		stoppableItems: items,
	}
	shutDownChan := make(chan os.Signal, 1)
	h.signalHandler = signal.NewSignalHandler(shutDownChan, &h)

	// WHEN
	stopable1.EXPECT().String().Return("stopable1")
	stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR"))
	stopable2.EXPECT().String().Return("stopable2")
	stopable2.EXPECT().Stop().Return(nil)

	start := time.Now()
	go func() {
		time.Sleep(time.Second * 1)
		shutDownChan <- signalMock{}
	}()

	time.Sleep(time.Millisecond * 100)
	h.WaitUntilSignal()

	// THEN
	assert.WithinDuration(t, start, time.Now(), time.Second*340)
}

func Test_RegisterFront(t *testing.T) {
	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stopable1 := NewMockStoppable(mockCtrl)
	stopable2 := NewMockStoppable(mockCtrl)
	h := ShutdownHandler{
		logger:         zerolog.Nop(),
		stoppableItems: &stop.OrderedStoppableList{},
	}
	shutDownChan := make(chan os.Signal, 1)
	h.signalHandler = signal.NewSignalHandler(shutDownChan, &h)

	// WHEN
	gomock.InOrder(
		stopable2.EXPECT().String().Return("stopable2"),
		stopable2.EXPECT().Stop().Return(nil),
		stopable1.EXPECT().String().Return("stopable1"),
		stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR")),
	)

	start := time.Now()
	go func() {
		time.Sleep(time.Second * 1)
		shutDownChan <- signalMock{}
	}()

	h.Register(stopable1)
	h.Register(stopable2)

	h.WaitUntilSignal()
	time.Sleep(time.Millisecond * 100)

	// THEN
	assert.WithinDuration(t, start, time.Now(), time.Millisecond*1200)
}

func Test_RegisterBack(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stopable1 := NewMockStoppable(mockCtrl)
	stopable2 := NewMockStoppable(mockCtrl)
	h := ShutdownHandler{
		logger:         zerolog.Nop(),
		stoppableItems: &stop.OrderedStoppableList{},
	}
	shutDownChan := make(chan os.Signal, 1)
	h.signalHandler = signal.NewSignalHandler(shutDownChan, &h)

	// WHEN
	gomock.InOrder(
		stopable1.EXPECT().String().Return("stopable1"),
		stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR")),
		stopable2.EXPECT().String().Return("stopable2"),
		stopable2.EXPECT().Stop().Return(nil),
	)

	start := time.Now()
	go func() {
		time.Sleep(time.Second * 1)
		shutDownChan <- signalMock{}
	}()

	h.Register(stopable1, false)
	h.Register(stopable2, false)

	h.WaitUntilSignal()
	time.Sleep(time.Millisecond * 100)

	// THEN
	assert.WithinDuration(t, start, time.Now(), time.Millisecond*1200)
}
