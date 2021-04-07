package common

import (
	"echo-stripe/response"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	respCode := 500
	resp := response.BasicResponse{}
	resp.Success = false
	resp.Message = ""

	sendErrorResponse := func() {
		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(respCode)
			} else {
				err = c.JSON(respCode, resp)
			}
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		// We have validation errors, sending back a 400.
		respCode = 400
		resp.Message = "Validation errors"
		resp.Errors = map[string]string{}

		for _, ve := range validationErrors {
			resp.Errors[ve.Field()] = "Contains unexpected value"
		}

		sendErrorResponse()
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		respCode = he.Code
		resp.Message = fmt.Sprintf("%v", he.Message)

		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)

			resp.Message += " - ErrorMessageInternal: " + he.Internal.Error()
		}
	} else {
		resp.Message = http.StatusText(respCode)
	}
	//else if MainCfg.Development {
	//	resp.Message = err.Error()
	//}

	if respCode == 500 {
		// 4 KB stack.
		stack := make([]byte, 4<<10)
		length := runtime.Stack(stack, false)
		fmt.Printf("[RECOVER From Exception]: %v %s\n", err, stack[:length])
	}

	c.Logger().Debug(err)
	sendErrorResponse()
}
