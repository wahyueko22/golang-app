package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/test/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostAuthOK(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	h.C.Orm.DB().DB = db

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "admin", "$2a$10$WUwK.b4F6BoXjBoq1ORpTONnXwrnoyA2EA7BfS9iNNEJRmkg8oGXq")

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	body := strings.NewReader(`{"login":"admin", "password":"admin"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}

func TestPostAuthFailCredentials(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	h.C.Orm.DB().DB = db

	rows := sqlmock.NewRows([]string{"id", "email", "password"})

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	body := strings.NewReader(`{"login":"admin@example.com", "password":"admin"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestPostAuthFailPassword(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	h.C.Orm.DB().DB = db

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "admin", "$2a$10$WUwK.b4F6BoXjBoq1ORpTONnXwrnoyA2EA7BfS9iNNEJRmkg8oGXq")

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	body := strings.NewReader(`{"login":"admin", "password":"adminaaa"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestPostAuthFailBodyFail(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())

	body := strings.NewReader(`or 1=1; select * from users;`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestPostAuthFailBodyLogin(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())

	body := strings.NewReader(`{"email":"admin", "password":"admin"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestPostAuthFailBodyPassword(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())

	body := strings.NewReader(`{"login":"admin", "passwd":"admin"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}

func TestPostAuthFailDbFind(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	h.C.Orm.DB().DB = db

	mock.ExpectBegin().WillReturnError(errors.New("mocked error"))
	mock.ExpectRollback()

	body := strings.NewReader(`{"login":"admin", "password":"admin"}`)
	resp, c := testutils.SetTestEnv(echo.POST, "/auth", body)

	if assert.NoError(t, h.PostAuth(c)) {
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	}
}
