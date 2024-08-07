package shutdown

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

type signalMock struct {
}

func (s signalMock) Signal() {
}
func (s signalMock) String() string {
	return ""
}

func Test_ShutdownHandler(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	var stopables []Stopable
	stopable1 := NewMockStopable(mockCtrl)
	stopables = append(stopables, stopable1)
	stopable2 := NewMockStopable(mockCtrl)
	stopables = append(stopables, stopable2)
	var logger zerolog.Logger
	h := Handler{
		orderedStopables:  stopables,
		isShutdownPending: atomic.NewBool(false),
		mux:               &sync.Mutex{},
	}
	shutDownChan := make(chan os.Signal, 1)

	// WHEN
	stopable1.EXPECT().String().Return("stopable1")
	stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR"))
	stopable2.EXPECT().String().Return("stopable2")
	stopable2.EXPECT().Stop().Return(nil)
	h.wg.Add(1)
	go h.shutdownHandler(shutDownChan, logger)

	start := time.Now()
	go func() {
		time.Sleep(time.Second * 1)
		shutDownChan <- signalMock{}
	}()

	h.WaitUntilSignal()
	time.Sleep(time.Millisecond * 100)

	// THEN
	assert.WithinDuration(t, start, time.Now(), time.Millisecond*1200)
}

func Test_RegisterFront(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	stopable1 := NewMockStopable(mockCtrl)
	stopable2 := NewMockStopable(mockCtrl)
	var logger zerolog.Logger
	h := Handler{
		orderedStopables:  make([]Stopable, 0),
		isShutdownPending: atomic.NewBool(false),
		mux:               &sync.Mutex{},
	}
	shutDownChan := make(chan os.Signal, 1)

	// WHEN
	gomock.InOrder(
		stopable2.EXPECT().String().Return("stopable2"),
		stopable2.EXPECT().Stop().Return(nil),
		stopable1.EXPECT().String().Return("stopable1"),
		stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR")),
	)
	h.wg.Add(1)
	go h.shutdownHandler(shutDownChan, logger)

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
	stopable1 := NewMockStopable(mockCtrl)
	stopable2 := NewMockStopable(mockCtrl)
	var logger zerolog.Logger
	h := Handler{
		orderedStopables:  make([]Stopable, 0),
		isShutdownPending: atomic.NewBool(false),
		mux:               &sync.Mutex{},
	}
	shutDownChan := make(chan os.Signal, 1)

	// WHEN
	gomock.InOrder(
		stopable1.EXPECT().String().Return("stopable1"),
		stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR")),
		stopable2.EXPECT().String().Return("stopable2"),
		stopable2.EXPECT().Stop().Return(nil),
	)
	h.wg.Add(1)
	go h.shutdownHandler(shutDownChan, logger)

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
