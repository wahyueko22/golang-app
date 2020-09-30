package main

import (
	"io"
	"net/http"
	"os"

	"echo-crud/auth"

	"github.com/labstack/echo/v4"
	myMiddleware "github.com/labstack/echo/v4/middleware"
)

var (
	LoginSuccessCookieVal = "LoginSuccess"
	UserIdExample         = "20022012"
	SecretKeyExample      = "mLmHu8f1IxFo4dWurBG3jEf1Ex0wDZvvwND6eFmcaX"
	SigningMethodExample  = "HS512"
)

func main() {
	e := echo.New()

	// Custom Middleware
	e.Use(CustomMiddleWareForServerHeader)

	jwtGroup := e.Group("/jwt")

	// JWT Middleware
	jwtGroup.Use(myMiddleware.JWTWithConfig(myMiddleware.JWTConfig{
		SigningMethod: SigningMethodExample,
		SigningKey:    []byte(SecretKeyExample),
	}))

	// ROUTING `jwtGroup`
	jwtGroup.GET("/main", mainJwt)

	e.Logger.Info("masukkk")
	// ROUTING `E`
	e.GET("/login", auth.Login)
	e.GET("/", welcome)
	e.GET("/users/:id", getUser)
	e.POST("/users/form", saveUserByForm)
	e.POST("/users/json", saveUserByJson)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.POST("/upload", fileUpload)

	// SERVER HOST
	e.Logger.Fatal(e.Start(":1323"))
}

func mainJwt(c echo.Context) error {
	return c.String(http.StatusOK, "SUCCESS: you are on the top secret jwt page aaaaa!")
}

func CustomMiddleWareForServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "ServerFirman/1.0")
		c.Response().Header().Set("My-Custom-Header", "ThisHaveNoMeaning")
		return next(c)
	}
}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	// return c.String(http.StatusOK, "id: "+id+", team: "+team+", member: "+member)
	return c.JSON(http.StatusOK, map[string]string{
		"id":     id,
		"team":   team,
		"member": member,
	})
}

func saveUserByForm(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name: "+name+", email: "+email)
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, "ID: "+id+", successfully updated.")
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusNoContent, id)
}

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func saveUserByJson(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func fileUpload(c echo.Context) error {
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("upload-files/" + avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Your file `"+avatar.Filename+"` successfully uploaded")
}
