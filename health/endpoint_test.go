package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CheckEvaluationResultToResponse(t *testing.T) {

	// GIVEN  no check
	at := time.Now()
	_29SecAfter := at.Add(time.Second * 29)
	cer := checkEvaluationResult{at: at}
	timeout := time.Second * 30

	// WHEN
	status, response := checkEvaluationResultToResponse(cer, _29SecAfter, timeout)

	// THEN
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "healthy", response.Status)
	assert.Len(t, response.Checks, 0)
	assert.Equal(t, at, response.At)

	// GIVEN  multiple checks
	healthyness := make(map[string]error)
	healthyness["check1"] = fmt.Errorf("No connection")
	healthyness["check2"] = fmt.Errorf("Timeout")
	healthyness["check3"] = nil

	cer = checkEvaluationResult{at: at, checkHealthyness: healthyness, numErrors: 2}

	// WHEN
	status, response = checkEvaluationResultToResponse(cer, _29SecAfter, timeout)

	// THEN
	assert.Equal(t, http.StatusServiceUnavailable, status)
	assert.Equal(t, "unhealthy", response.Status)
	assert.Len(t, response.Checks, 3)

	checkByName := make(map[string]check)
	checkByName[response.Checks[0].Name] = response.Checks[0]
	checkByName[response.Checks[1].Name] = response.Checks[1]
	checkByName[response.Checks[2].Name] = response.Checks[2]
	assert.Equal(t, "No connection", checkByName["check1"].Error)
	assert.Equal(t, "Timeout", checkByName["check2"].Error)
	assert.Empty(t, checkByName["check3"].Error)
}

func Test_CheckEvaluationResultToResponseShouldBeUnhealthyIfTooOld(t *testing.T) {

	// GIVEN
	at := time.Now()
	_31SecAfter := at.Add(time.Second * 31)
	cer := checkEvaluationResult{at: at}
	timeout := time.Second * 30

	// WHEN
	status, response := checkEvaluationResultToResponse(cer, _31SecAfter, timeout)

	// THEN
	assert.Equal(t, http.StatusServiceUnavailable, status)
	assert.Equal(t, "unhealthy", response.Status)
	assert.Len(t, response.Checks, 0)
	assert.Equal(t, at, response.At)
}

func Test_HealthEndpoint(t *testing.T) {

	// GIVEN
	registry := NewRegistry()
	monitor, err := NewMonitor(&registry)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	healthyness := make(map[string]error)
	healthyness["check1"] = nil
	healthyness["check2"] = fmt.Errorf("Timeout")
	monitor.latestCheckResult.Store(checkEvaluationResult{
		at:               time.Now(),
		numErrors:        2,
		checkHealthyness: healthyness,
	})

	// WHEN
	monitor.Health(w, req)

	// THEN
	resp := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	defer resp.Body.Close()

	respHealth := response{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&respHealth)
	require.NoError(t, err)
	assert.WithinDuration(t, respHealth.At, time.Now(), time.Millisecond*10)
	assert.Equal(t, "unhealthy", respHealth.Status)
	assert.Len(t, respHealth.Checks, 2)

	checkByName := make(map[string]check)
	checkByName[respHealth.Checks[0].Name] = respHealth.Checks[0]
	checkByName[respHealth.Checks[1].Name] = respHealth.Checks[1]
	assert.Equal(t, "healthy", checkByName["check1"].Status)
	assert.Empty(t, checkByName["check1"].Error)
	assert.Equal(t, "unhealthy", checkByName["check2"].Status)
	assert.Equal(t, "Timeout", checkByName["check2"].Error)
}
