package common

import (
	"echo-crud/model"

	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	User *model.User
}
