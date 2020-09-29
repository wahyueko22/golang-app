package version

import (
	"net/http"
	"testing"

	"github.com/corvinusz/echo-xorm/test/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	resp, c := testutils.SetTestEnv(echo.GET, "/version", nil)
	h := NewHandler(testutils.SetTestAppContext())

	var (
		versionJSON = `"version":"0.1.0-dev"`
		resultJSON  = `"result":"OK"`
	)

	if assert.NoError(t, h.GetVersion(c)) {
		assert.Equal(t, http.StatusOK, resp.Code)
		respBody := resp.Body.String()
		assert.Contains(t, respBody, versionJSON)
		assert.Contains(t, respBody, resultJSON)
	}
}
