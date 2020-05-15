package controller

// User 控制器
type User struct {
}

// GetUserBy 根据 username 查找用户
func (u *User) GetUserBy(username string) string {
	return "wapapi get: " + username
}
