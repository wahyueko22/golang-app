package auth

import (
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

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "wahyu" && password == "password" {
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
