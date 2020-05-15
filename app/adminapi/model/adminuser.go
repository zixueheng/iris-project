package model

import (
	"iris-project/app"
	"iris-project/app/adminapi/validate"
	"iris-project/config"
	"iris-project/global"
	"iris-project/lib/util"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

// AdminUser 模型
type AdminUser struct {
	gorm.Model
	Username            string `gorm:"type:varchar(100);unique_index;not null"`
	Password            string `gorm:"type:varchar(100);not null"`
	Token               string
	TokenExpired        time.Time
	RefreshToken        string
	RefreshTokenExpired time.Time
	Role                Role
	RoleID              uint
	Phone               string `gorm:"type:char(11);not null"`
	Status              int8   `gorm:"type:tinyint(1);default:1"`
}

// FindOne 查找一个
func FindOne(obj interface{}, where map[string]interface{}, preload string) bool {
	tx := global.Db.Where(where)

	// 这里 Preload() 不起作用
	if len(preload) != 0 {
		tx.Preload(preload)
	}

	if err := tx.First(obj).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// CheckLogin 登录检查
func (au *AdminUser) CheckLogin(loginInfo *validate.LoginRequest) (interface{}, bool, app.Code) {
	if !au.GetAdminUserByUsername(loginInfo.Username) {
		return nil, false, app.CodeUserNotFound
	} else {
		if bcrypt.Match(loginInfo.Password, au.Password) {
			var now = time.Now()
			var tokenExpired = now.Add(time.Minute * time.Duration(10))
			var refreshTokenExpired = now.Add(time.Minute * time.Duration(30))

			// 获取一个 Token，参数一：签名方法、参数二：要保存的数据
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"admin_user_id": util.ParseString(int(au.ID)),
				"exp":           util.TimeFormat(tokenExpired, ""),
				"iat":           util.TimeFormat(now, ""),
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, _ := token.SignedString([]byte(config.App.Jwtsecret))

			refreshToken := util.GetRandomString(32)

			au.Token = tokenString
			au.TokenExpired = tokenExpired
			au.RefreshToken = refreshToken
			au.RefreshTokenExpired = refreshTokenExpired
			global.Db.Save(au) // 保存 token 等信息

			// TODO 这里 token等信息最好保存到 Redis 里面

			response := struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
			}{
				Token:        tokenString,
				RefreshToken: refreshToken,
			}

			return response, true, app.CodeUserLoginSucceed

		}
		return nil, false, app.CodeUserPasswordError
	}

}

// GetAdminUserByUsername 根据用户名获取管理员
func (au *AdminUser) GetAdminUserByUsername(username string) bool {
	if err := global.Db.Where("username=?", username).Preload("Role").First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

// CreateAdminUser 创建管理员
func (au *AdminUser) CreateAdminUser() error {
	if global.Db.NewRecord(au) {
		if err := global.Db.Create(&au).Error; err != nil {
			return err
		}
	}

	return nil
}
