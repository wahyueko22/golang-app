package logger

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// HTTPLogger returns a middleware that logs HTTP requests.
func HTTPLogger(l Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			var logData = map[string]interface{}{
				"Remote":      c.RealIP(),
				"Method":      req.Method,
				"Proto":       req.Proto,
				"URL":         req.URL.String(),
				"Status":      strconv.Itoa(c.Response().Status),
				"ElapsedTime": stop.Sub(start).String(),
				// "Data":    req.GetBody(),
				// "AuthData":    req.Header.Get(""),
			}
			// pack data and send to logger
			l.Info("http", fmt.Sprint(logData))
			return nil
		}
	}
}
