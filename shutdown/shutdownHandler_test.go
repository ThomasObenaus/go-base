package shutdown

import (
	"fmt"
	"os"
	"testing"
	"time"

	mock_shutdown "github.com/ThomasObenaus/go-base/test/mocks/shutdown"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
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
	stopable1 := mock_shutdown.NewMockStopable(mockCtrl)
	stopables = append(stopables, stopable1)
	stopable2 := mock_shutdown.NewMockStopable(mockCtrl)
	stopables = append(stopables, stopable2)
	var logger zerolog.Logger
	h := Handler{}
	shutDownChan := make(chan os.Signal, 1)

	// WHEN
	stopable1.EXPECT().Name().Return("stopable1")
	stopable1.EXPECT().Stop().Return(fmt.Errorf("ERROR"))
	stopable2.EXPECT().Name().Return("stopable2")
	stopable2.EXPECT().Stop().Return(nil)
	go h.shutdownHandler(shutDownChan, stopables, logger)
	shutDownChan <- signalMock{}
	time.Sleep(time.Second * 1)

	// THEN
}
