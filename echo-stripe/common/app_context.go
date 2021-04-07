package common

import (
	"echo-stripe/model"

	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	User *model.User
}
