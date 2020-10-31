package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	UserId        = "USER_123456"
	SecretKey     = "secret"
	SigningMethod = "HS512"
)

type LoginForm struct {
	UserName string `json:"username" form:"username" query:"username"  validate:"required"`
	Password string `json:"password" form:"password" query:"password"  validate:"required"`
}

func Login(c echo.Context) error {
	//username := c.QueryParam("username")
	//password := c.QueryParam("password")
	//currContext := c.(*common.AppContext)
	//u := &user{}
	form := new(LoginForm)
	//form := &LoginForm{}
	fmt.Println(" masukkk  loginnn : ")

	if err := c.Bind(&form); err != nil {
		return err
	}

	if err := c.Validate(form); err != nil {
		//log.Print(err)
		return err
	}

	if form.UserName == "wahyu" && form.Password == "password" {
		// TODO: create jwt token
		token, err := createJwtToken()
		if err != nil {
			log.Println("Error when creating JWT token", err)
			return c.String(http.StatusInternalServerError, "ERROR: something went wrong while creating JWT token!")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You ware logged in!",
			"token":   token,
		})
	}

	return c.String(http.StatusUnauthorized, "WARNING: Make sure your account is coorect!")
}

type JwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func createJwtToken() (string, error) {
	claims := JwtCustomClaims{
		"wahyu",
		jwt.StandardClaims{
			Id:        UserId,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// we hash the jwt claims
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
