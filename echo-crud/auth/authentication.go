package auth

import (
	"log"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var (
	UserIdExample    = "20022012"
	SecretKeyExample = "mLmHu8f1IxFo4dWurBG3jEf1Ex0wDZvvwND6eFmcaX"
)

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "wahyu" && password == "passwprd" {
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

type JwtClaims struct {
	Name string `json:"name"`
	jwtGo.StandardClaims
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"Firman",
		jwtGo.StandardClaims{
			Id:        UserIdExample,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// we hash the jwt claims
	rawToken := jwtGo.NewWithClaims(jwtGo.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte(SecretKeyExample))
	if err != nil {
		return "", err
	}

	return token, nil
}
