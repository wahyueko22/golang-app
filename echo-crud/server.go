package main

import (
	"fmt"
	"os"
	"strings"

	"echo-crud/auth"
	"echo-crud/common"
	"echo-crud/model"
	"echo-crud/services"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	appMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go run server.go
func main() {
	e := echo.New()
	//setting log level
	e.Logger.SetLevel(log.DEBUG)
	// Custom Middleware
	e.Use(customServerHeader)

	e.HTTPErrorHandler = common.CustomHTTPErrorHandler

	//field validation https://godoc.org/github.com/go-playground/validator
	custValidator := &common.CustomValidator{Validator: validator.New()}
	err := custValidator.Init()
	if err != nil {
		fmt.Printf("Failed to init validator: %v\n", err)
		os.Exit(1)
	}
	e.Validator = custValidator

	myFile, err := os.Create("logApp.json")
	if err != nil {
		panic(err)
	}
	e.Use(appMiddleware.LoggerWithConfig(appMiddleware.LoggerConfig{
		//Format:  "method=${method}, uri=${uri}, status=${status} latency=${latency_human} in:${bytes_in} out:${bytes_out}\n",
		//Skipper: DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		//Output: os.Stdout,
		Output: myFile,
	}))

	//allow CORS
	e.Use(appMiddleware.CORS())

	//https://echo.labstack.com/guide/context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//reference data
			appContext := &common.AppContext{}
			appContext.Context = c
			return h(appContext)
		}
	})

	// JWT Middleware. middle is interceptor
	e.Use(appMiddleware.JWTWithConfig(appMiddleware.JWTConfig{
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
	fmt.Println(" masukkk  println : ")
	// ROUTING `E`
	e.POST("/login", auth.Login)
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
		c.Response().Header().Set(echo.HeaderServer, "GolangApp/1.0")
		c.Response().Header().Set("My-Custom-Header", "CustomAdditionalHeader")
		return next(c)
	}
}
