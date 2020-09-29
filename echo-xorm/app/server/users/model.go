package users

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/corvinusz/echo-xorm/pkg/errors"
	"github.com/corvinusz/echo-xorm/pkg/utils"

	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity (here are DB definitions)
type User struct {
	ID            uint64         `xorm:"'id' pk autoincr unique notnull" json:"id"`
	Email         string         `xorm:"'email' text index not null unique" json:"email"`
	DisplayName   string         `xorm:"'display_name' text" json:"displayName"`
	Password      string         `xorm:"'password' text not null" json:"-"`
	PasswordEtime uint64         `xorm:"'password_etime'" json:"passwordEtime"`
	PasswordURL   sql.NullString `xorm:"'password_url' text unique" json:"passwordUrl"`
	Created       uint64         `xorm:"created" json:"created"`
	Updated       uint64         `xorm:"updated" json:"updated"`
}

// TableName used by xorm to set table name for entity
func (u *User) TableName() string {
	return "users"
}

// NewUser creates user from request body
// returns *User with data from body
// returns nil if error occurred
func NewUser(b *PostBody) *User {
	u := &User{
		Email:         b.Email,
		DisplayName:   b.DisplayName,
		Password:      b.Password,
		PasswordEtime: b.PasswordEtime,
	}
	// passwordURL
	switch {
	case b.PasswordURL == nil:
		u.PasswordURL = sql.NullString{
			Valid: false,
		}
	case *b.PasswordURL == "":
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: "",
		}
	default:
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: *b.PasswordURL,
		}

	}
	return u
}

// FindAll users in database
func FindAll(orm *xorm.Engine) ([]User, error) {
	users := []User{}
	err := orm.Find(&users)
	if err != nil {
		return nil, errors.New("database error; " + err.Error())
	}
	return users, nil
}

// FindOne finds in database first user with matched struct fields
func (u *User) FindOne(orm *xorm.Engine) error {
	found, err := orm.Get(u)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}
	return nil
}

// Save user to database
func (u *User) Save(orm *xorm.Engine) error {
	// begin transaction
	tx, err := utils.BeginTransaction(orm)
	if err != nil {
		return err
	}
	defer tx.Close()
	// validate data
	err = u.validateDataToSave(tx)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	// save data to storage
	affected, err := tx.InsertOne(u)
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	if affected == 0 {
		return errors.New("database error; db refused to insert")
	}
	// commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}

	return nil
}

// Update user in database
func (u *User) Update(orm *xorm.Engine) error {
	// begin transaction
	tx, err := utils.BeginTransaction(orm)
	if err != nil {
		return err
	}
	defer tx.Close()
	// get old user data (and check if user exists)
	old := &User{ID: u.ID}
	found, err := tx.Get(old)
	if err != nil {
		err = errors.NewWithPrefix(err, "database error")
		return utils.RollbackTransaction(tx, err)
	}
	if !found {
		err = errors.NewWithCode(http.StatusNotFound, "user not found")
		return utils.RollbackTransaction(tx, err)
	}
	// set data to update
	err = u.setDataToUpdate(old)
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	// validate data to update
	err = u.validateDataToUpdate(tx)
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	// update
	affected, err := tx.ID(u.ID).Update(u)
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	if affected == 0 {
		err = errors.New("database error; db refused to update")
		return utils.RollbackTransaction(tx, err)
	}
	// commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}

	return nil
}

// Delete user from database
func (u *User) Delete(orm *xorm.Engine) error {
	// begin transaction
	tx, err := utils.BeginTransaction(orm)
	if err != nil {
		return err
	}
	defer tx.Close()
	// check if user exists
	old := &User{ID: u.ID}
	found, err := tx.Get(old)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}
	// delete from storage
	affected, err := tx.ID(u.ID).Delete(&User{})
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	if affected == 0 {
		return errors.New("db refused to delete")
	}
	// commit
	err = tx.Commit()
	if err != nil {
		return utils.RollbackTransaction(tx, err)
	}
	return nil
}

//------------------------------------------------------------------------------
func (u *User) setDataToUpdate(old *User) error {
	// email
	if u.Email == "" {
		u.Email = old.Email
	}
	// displayName
	if u.DisplayName == "" {
		u.DisplayName = old.DisplayName
	}
	// passwordEtime
	if u.PasswordEtime == 0 {
		u.PasswordEtime = old.PasswordEtime
	}
	// password
	if u.Password == "" {
		u.Password = old.Password
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.NewWithPrefix(err, "password hash generation error")
		}
		u.Password = string(hash[:8])
	}
	// created/updated
	u.Created = old.Created
	u.Updated = uint64(time.Now().UTC().Unix())
	return nil
}

func (u *User) validateDataToSave(tx *xorm.Session) error {
	// email uniqueness
	affected, err := tx.Where("email = ?", u.Email).Count(&User{})
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected != 0 {
		return errors.NewWithCode(http.StatusConflict, "email already exists")
	}

	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewWithPrefix(err, "generate hash from password error")
	}
	u.Password = string(hash)
	// created/updated
	u.Created = uint64(time.Now().UTC().Unix())
	u.Updated = u.Created
	return nil
}

func (u *User) validateDataToUpdate(tx *xorm.Session) error {
	// email uniqueness
	affected, err := tx.Where("email = ?", u.Email).
		And("id != ?", u.ID).Count(&User{})
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected != 0 {
		return errors.NewWithCode(http.StatusConflict, "email already exists")
	}
	return nil
}
