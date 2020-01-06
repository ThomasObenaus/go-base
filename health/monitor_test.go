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
	registry := NewRegistry()

	// WHEN
	monitor, err := NewMonitor(&registry)

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, monitor)
	assert.NotNil(t, monitor.registry)
}

func Test_NewMonitorShouldFail(t *testing.T) {

	// GIVEN

	// WHEN
	monitor, err := NewMonitor(nil)

	// THEN
	assert.Error(t, err)
	assert.Nil(t, monitor)
}

func Test_EvaluateChecks(t *testing.T) {

	// GIVEN
	registry := NewRegistry()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nameCheck1 := "check1-healthy"
	check1 := mock_health.NewMockCheck(mockCtrl)
	check1.EXPECT().Name().Return(nameCheck1).Times(2)
	check1.EXPECT().IsHealthy().Return(nil)
	err := registry.Register(check1)
	require.NoError(t, err)

	nameCheck2 := "check2-unhealthy"
	errCheck2 := fmt.Errorf("could not connect")
	check2 := mock_health.NewMockCheck(mockCtrl)
	check2.EXPECT().Name().Return(nameCheck2).Times(2)
	check2.EXPECT().IsHealthy().Return(errCheck2)
	err = registry.Register(check2)
	require.NoError(t, err)

	monitor, err := NewMonitor(&registry)
	require.NoError(t, err)
	require.NotNil(t, monitor)

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

func ExampleNewMonitor() {
	registry := NewRegistry()
	monitor, _ := NewMonitor(&registry)
	monitor.Start()

	// register the endpoint at the router/ server of your choice
	http.HandleFunc("/health", monitor.Health)
	http.ListenAndServe(":8080", nil)

	monitor.Stop()
}
