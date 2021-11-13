package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type response struct {
	At     time.Time `json:"at,omitempty"`
	Status string    `json:"status,omitempty"`
	Checks []check   `json:"checks"`
}

type check struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func checkEvaluationResultToResponse(cer checkEvaluationResult, now time.Time, checkEvaluationTimeout time.Duration) (int, response) {

	response := response{
		At:     cer.at,
		Status: "healthy",
	}

	httpStatusCode := http.StatusOK

	// switch to unhealthy in case there is at least one error
	// or the last evaluation was too long ago
	if cer.numErrors > 0 || now.Sub(cer.at) >= checkEvaluationTimeout {
		httpStatusCode = http.StatusServiceUnavailable
		response.Status = "unhealthy"
	}

	// add the result from the checks
	var checks []check

	for checkName, err := range cer.checkHealthyness {

		status := "healthy"
		errMsg := ""
		if err != nil {
			status = "unhealthy"
			errMsg = err.Error()
		}

		checks = append(checks, check{
			Name:   checkName,
			Status: status,
			Error:  errMsg,
		})
	}

	response.Checks = checks
	return httpStatusCode, response
}

// Health is the health endpoint
func (m *Monitor) Health(w http.ResponseWriter, r *http.Request) {
	m.logger.Debug().Msg("Health endpoint called")

	latestResult := m.latestCheckResult.Load().(checkEvaluationResult)
	code, response := checkEvaluationResultToResponse(latestResult, time.Now(), m.checkEvaluationTimeout)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
