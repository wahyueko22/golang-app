package common

import "github.com/labstack/echo/v4"

func CheckError(err error) {
	if err == nil {
		return
	}
	//	error.new("test")
	if _, ok := err.(*echo.HTTPError); ok {
		panic(err)
	}

	AbortErr(err)
}

func Abort(statusCode int) {
	panic(echo.NewHTTPError(statusCode))
}

func AbortErr(err error) {
	ret := echo.NewHTTPError(500)
	ret.SetInternal(err)
	panic(ret)
}

func AbortWithMessage(statusCode int, msg string) {
	ret := echo.NewHTTPError(statusCode)
	ret.Message = msg
	panic(ret)
}
