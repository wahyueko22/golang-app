package users

import (
	"net/http"
	"testing"

	"github.com/corvinusz/echo-xorm/test/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersOK(t *testing.T) {
	h := NewHandler(testutils.SetTestAppContext())
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	h.C.Orm.DB().DB = db

	rows := sqlmock.NewRows([]string{"id", "email", "display_name"}).
		AddRow(1, "admin", "admin")

	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	resp, c := testutils.SetTestEnv(echo.GET, "/operators", nil)

	const (
		expectedEmail = `"email":"admin"`
		expectedName  = `"displayName":"admin"`
	)

	if assert.NoError(t, h.GetAllUsers(c)) {
		assert.Equal(t, http.StatusOK, resp.Code)
		respBody := resp.Body.String()
		assert.Contains(t, respBody, expectedEmail)
		assert.Contains(t, respBody, expectedName)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
