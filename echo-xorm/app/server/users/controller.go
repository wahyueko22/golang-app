package users

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/corvinusz/echo-xorm/app/ctx"
	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/pkg/utils"

	"github.com/labstack/echo/v4"
)

var reEmail = regexp.MustCompile(ctx.EmailValidation)

// PostBody represents payload data format
type PostBody struct {
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	DisplayName   string  `json:"displayName"`
	PasswordURL   *string `json:"passwordUrl"`
	PasswordEtime uint64  `json:"passwordEtime"`
}

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

func NewHandler(c *ctx.Context) *Handler {
	return &Handler{C: c}
}

// GetAllUsers is a GET /users handler
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := FindAll(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser is a GET /users/{id} handler
func (h *Handler) GetUser(c echo.Context) error {
	var (
		user User
		err  error
	)

	user.ID, err = strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request paramer read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}

	err = user.FindOne(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser is a POST /users handler
func (h *Handler) PostUser(c echo.Context) error {
	var body PostBody
	err := c.Bind(&body)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// validate body
	err = validatePostBody(&body)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body validate error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// create
	user := NewUser(&body)
	// save
	err = user.Save(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusCreated, user)
}

// PutUser is a PUT /users/{id} handler
func (h *Handler) PutUser(c echo.Context) error {
	var body PostBody
	// parse id
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request parameter read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// parse request body
	if err = c.Bind(&body); err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// validate body
	err = validatePutBody(&body)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request body validate error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	// construct user
	user := NewUser(&body)
	user.ID = id
	// update
	err = user.Update(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser is a DELETE /users/{id} handler
func (h *Handler) DeleteUser(c echo.Context) error {
	var user User

	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		err = errors.NewWithCode(http.StatusBadRequest, "request paramer read error; "+err.Error())
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}

	user.ID = id
	// delete
	err = user.Delete(h.C.Orm)
	if err != nil {
		h.C.Logger.Error(utils.GetEvent(c), err.Error())
		return c.String(errors.Decompose(err))
	}
	return c.NoContent(http.StatusOK)
}

func validatePostBody(b *PostBody) error {
	if !reEmail.MatchString(b.Email) {
		if b.Email != "admin" {
			return errors.New("invalid email")
		}
	}
	if len(b.Password) < 6 || len(b.Password) > 100 {
		return errors.New("invalid password")
	}
	return nil
}

func validatePutBody(b *PostBody) error {
	if (b.Email != "") && !reEmail.MatchString(b.Email) {
		if b.Email != "admin" {
			return errors.New("invalid email")
		}
	}
	if (b.Password != "") && (len(b.Password) < 6 || len(b.Password) > 100) {
		return errors.New("invalid password")
	}
	return nil
}
