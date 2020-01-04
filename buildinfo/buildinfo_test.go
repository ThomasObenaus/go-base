package buildinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Print(t *testing.T) {

	var builder strings.Builder
	testPrintFun := func(format string, a ...interface{}) (int, error) {
		builder.WriteString(fmt.Sprintf(format, a))
		return 0, nil
	}

	// GIVEN
	bi := BuildInfo{Version: "v1.0.0", BuildTime: "2020-01-04", Revision: "ABCED", Branch: "featurebranch"}

	// WHEN
	bi.Print(testPrintFun)

	// THEN
	str := builder.String()
	assert.Contains(t, str, "v1.0.0")
	assert.Contains(t, str, "2020-01-04")
	assert.Contains(t, str, "ABCED")
	assert.Contains(t, str, "featurebranch")
}

func Test_PrintNoVersion(t *testing.T) {

	var builder strings.Builder
	testPrintFun := func(format string, a ...interface{}) (int, error) {
		builder.WriteString(fmt.Sprintf(format, a))
		return 0, nil
	}

	// GIVEN
	bi := BuildInfo{}

	// WHEN
	bi.Print(testPrintFun)

	// THEN
	str := builder.String()
	assert.Regexp(t, "Version:.*n/a", str)
	assert.Regexp(t, "Build-Time:.*n/a", str)
	assert.Regexp(t, "Revision:.*n/a", str)
	assert.Regexp(t, "on.*n/a", str)
}

func Test_BuildInfoEndpoint(t *testing.T) {

	bi := BuildInfo{Version: "v1.0.0", BuildTime: "2020-01-04", Revision: "ABCED", Branch: "featurebranch"}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	bi.BuildInfo(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	respBi := BuildInfo{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&respBi)
	require.NoError(t, err)
	assert.Equal(t, respBi.Version, bi.Version)
	assert.Equal(t, respBi.BuildTime, bi.BuildTime)
	assert.Equal(t, respBi.Revision, bi.Revision)
	assert.Equal(t, respBi.Branch, bi.Branch)
}
