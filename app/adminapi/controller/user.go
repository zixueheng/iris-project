package controller

import (
	"iris-project/app/model"

	"github.com/kataras/iris/v12"
)

// User 控制器
type User struct {
	Ctx iris.Context // IRIS框架会自动注入 Context
}

// GetUserBy 根据 username 查找用户
func (u *User) GetUserBy(username string, age uint, birthday string) string {
	u.Ctx.Application().Logger().Info("adminapi get: " + username)

	user := model.User{Name: username, Age: age, Birthday: birthday}
	user.CreateUser()

	return "adminapi get: " + username
}
