package model

type User struct {
	Id    string `xorm:pk "id"`
	Name  string `xorm:pk "name"`
	Email string `xorm:pk "email"`
}
