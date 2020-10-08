package main

import (
	"fmt"
	"strings"

	"echo-crud/auth"
	"echo-crud/common"
	"echo-crud/model"
	"echo-crud/services"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	myMiddleware "github.com/labstack/echo/v4/middleware"
)

//go run server.go
func main() {
	e := echo.New()

	// Custom Middleware
	e.Use(customServerHeader)
	//https://echo.labstack.com/guide/context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//reference data
			appContext := &common.AppContext{}
			appContext.Context = c
			return h(appContext)
		}
	})

	// JWT Middleware
	e.Use(myMiddleware.JWTWithConfig(myMiddleware.JWTConfig{
		SigningMethod: auth.SigningMethod,
		SigningKey:    []byte(auth.SecretKey),
		Claims:        &auth.JwtCustomClaims{},
		Skipper: func(c echo.Context) bool {
			if strings.HasSuffix(c.Path(), "/login") {
				return true
			}
			return false
		},
		SuccessHandler: func(c echo.Context) {
			user, ok := c.Get("user").(*jwt.Token)
			if ok {
				claims := user.Claims.(*auth.JwtCustomClaims)
				//custom loader https://echo.labstack.com/guide/context
				appContext := c.(*common.AppContext)
				user := &model.User{}
				user.Name = claims.Name
				appContext.User = user
				fmt.Println(" Claim ID : ")
				fmt.Println(claims.Id)
				fmt.Println("Standart Claim ID: ")
				fmt.Println(claims.StandardClaims.Id)
			}
		},
	}))
	//	e.
	e.Logger.Info("masukkk")
	// ROUTING `E`
	e.GET("/login", auth.Login)
	e.GET("/", services.Welcome)
	e.GET("/main", services.MainJwt)
	e.GET("/users/:id", services.GetUser)
	e.POST("/users/form", services.SaveUserByForm)
	e.POST("/users/json", services.SaveUserByJson)
	e.PUT("/users/:id", services.UpdateUser)
	e.DELETE("/users/:id", services.DeleteUser)
	e.POST("/upload", services.FileUpload)

	// SERVER HOST
	e.Logger.Fatal(e.Start(":1323"))
}

func customServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "ServerFirman/1.0")
		c.Response().Header().Set("My-Custom-Header", "CustomAdditionalHeader")
		return next(c)
	}
}
