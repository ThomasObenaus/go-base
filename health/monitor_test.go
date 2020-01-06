package health

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	mock_health "github.com/ThomasObenaus/go-base/test/mocks/health"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewMonitor(t *testing.T) {

	// GIVEN

	// WHEN
	monitor, err := NewMonitor()

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, monitor)
}

func Test_EvaluateChecks(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nameCheck1 := "check1-healthy"
	check1 := mock_health.NewMockCheck(mockCtrl)
	check1.EXPECT().String().Return(nameCheck1).Times(2)
	check1.EXPECT().IsHealthy().Return(nil)

	nameCheck2 := "check2-unhealthy"
	errCheck2 := fmt.Errorf("could not connect")
	check2 := mock_health.NewMockCheck(mockCtrl)
	check2.EXPECT().String().Return(nameCheck2).Times(2)
	check2.EXPECT().IsHealthy().Return(errCheck2)

	monitor, err := NewMonitor()
	require.NoError(t, err)
	require.NotNil(t, monitor)
	err = monitor.Register(check1)
	require.NoError(t, err)
	err = monitor.Register(check2)
	require.NoError(t, err)

	// WHEN
	now := time.Now()
	checkResult := monitor.evaluateChecks(now)

	// THEN
	assert.Equal(t, now, checkResult.at)
	assert.Equal(t, uint(1), checkResult.numErrors)
	assert.Len(t, checkResult.checkHealthyness, 2)
	assert.Nil(t, checkResult.checkHealthyness[nameCheck1])
	assert.NotNil(t, checkResult.checkHealthyness[nameCheck2])
	assert.Equal(t, errCheck2, checkResult.checkHealthyness[nameCheck2])
}

func Test_ShouldRegister(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	check1 := mock_health.NewMockCheck(mockCtrl)
	monitor, err := NewMonitor()
	require.NoError(t, err)
	require.NotNil(t, monitor)

	// WHEN
	check1.EXPECT().String().Return("check1")
	err = monitor.Register(check1)

	// THEN
	assert.NoError(t, err)
	assert.Len(t, monitor.healthChecks, 1)
}

func Test_ShouldNotRegister(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	check1 := mock_health.NewMockCheck(mockCtrl)
	monitor, err := NewMonitor()
	require.NoError(t, err)
	require.NotNil(t, monitor)

	// WHEN
	check1.EXPECT().String().Return("")
	err = monitor.Register(check1)

	// THEN
	assert.Error(t, err)
	assert.Len(t, monitor.healthChecks, 0)

	// WHEN
	err = monitor.Register(nil)

	// THEN
	assert.Error(t, err)
	assert.Len(t, monitor.healthChecks, 0)
}

func TestRunJoinStop(t *testing.T) {

	// GIVEN
	monitor, err := NewMonitor()
	require.NotNil(t, monitor)
	require.NoError(t, err)

	monitor.Start()
	start := time.Now()
	monitor.Stop()
	monitor.Join()

	assert.WithinDuration(t, start.Add(time.Millisecond*500), time.Now(), time.Second*1)
}

func ExampleNewMonitor() {
	monitor, _ := NewMonitor()
	check, _ := NewSimpleCheck("my-check", func() error {
		// return nil if healthy
		// return the error if unhealthy
		return nil
	})

	monitor.Register(check)
	monitor.Start()

	// register the endpoint at the router/ server of your choice
	http.HandleFunc("/health", monitor.Health)
	http.ListenAndServe(":8080", nil)

	monitor.Stop()
}
