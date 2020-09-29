package testutils

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/pkg/logger"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func SetTestAppContext() *ctx.Context {
	testLogger := logger.NewStdLogger("unit-tests", "echo-xorm-test")
	testOrm, err := xorm.NewEngine("sqlite3", ":memory:") // fake database
	if err != nil {
		panic(err)
	}
	testConfig := ctx.Config{
		Version: "0.1.0-dev",
	}

	return &ctx.Context{
		Logger: testLogger,
		Orm:    testOrm,
		Config: &testConfig,
	}
}

func SetTestEnv(method, path string, body *strings.Reader) (rec *httptest.ResponseRecorder, c echo.Context) {
	e := echo.New()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	return
}
