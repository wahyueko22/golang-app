package auth

import (
	"net/http"
	"time"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/app/server/users"
	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/pkg/utils"

	jwt "github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Handler represents handlers for '/auth'
type Handler struct {
	C *ctx.Context
}

func NewHandler(c *ctx.Context) *Handler {
	return &Handler{C: c}
}

// PostBody represents payload data format
type PostBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Result represents payload response format
type Result struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

// PostAuth is handler for /auth
func (h *Handler) PostAuth(c echo.Context) error {
	var body PostBody

	err := c.Bind(&body)
	if err != nil {
		err = errors.NewWithPrefix(err, "request body parse error")
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(http.StatusUnauthorized, err.Error())
	}

	// find user
	user := users.User{Email: body.Email}
	err = user.FindOne(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.NoContent(http.StatusUnauthorized)
	}

	// validate user credentials
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		err = errors.NewWithPrefix(err, "compare hash and password")
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.NoContent(http.StatusUnauthorized)
	}

	// create a HMAC SHA256 signer
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = user.Email
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).UTC().Unix()
	claims["jti"] = user.ID

	t, err := token.SignedString(h.C.JWTSignKey)
	if err != nil {
		err = errors.NewWithPrefix(err, "token signing error")
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.NoContent(http.StatusUnauthorized)
	}

	resp := Result{
		Result: "OK",
		Token:  t,
	}
	return c.JSON(http.StatusOK, resp)
}
