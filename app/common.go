package app

import (
	"fmt"
	"iris-project/config"
	"iris-project/global"
	"iris-project/lib/util"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// APP公共函数

// CheckRequest 检查请求参数，返回错误提示
func CheckRequest(ctx iris.Context, obj interface{}) (errmsg string) {
	if err := ctx.ReadJSON(obj); err != nil {
		// p.Ctx.StatusCode(iris.StatusOK)
		// _, _ = ctx.JSON(APIData(false, nil, err.Error()))
		// return

		errmsg = err.Error()
	}
	// fmt.Printf("\n%+v\n", obj)
	err := global.Validate.Struct(obj)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(global.ValidateTrans) {
			// fmt.Println(e)
			if len(e) > 0 {
				// p.Ctx.StatusCode(iris.StatusOK)
				// _, _ = ctx.JSON(APIData(false, nil, e))
				// return

				errmsg += " " + e
				// break
			}
		}
	}
	return
}

// GenTokenAndRefreshToken 生成Token和RefreshToken
//
// key 键名，id 键值，
// tokenMinutes Token多少分钟后过期，refreshTokenMinutes 刷新Token多少分钟后过期
func GenTokenAndRefreshToken(key string, id int, tokenMinutes, refreshTokenMinutes int) (string, string) {
	var now = time.Now()
	var tokenExpired = now.Add(time.Minute * time.Duration(tokenMinutes))
	// var refreshTokenExpired = now.Add(time.Minute * time.Duration(10))

	// 获取一个 Token，参数一：签名方法、参数二：要保存的数据
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		key:   util.ParseString(id),
		"exp": util.TimeFormat(tokenExpired, ""),
		"iat": util.TimeFormat(now, ""),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(config.App.Jwtsecret))

	refreshToken := util.GetRandomString(64)

	// 保持刷新token到Redis中
	err := global.Redis.Set(config.App.Appname+":refresh_token_admin_"+util.ParseString(id), refreshToken, time.Minute*time.Duration(refreshTokenMinutes)).Err()
	if err != nil {
		panic(err)
	}
	return tokenString, refreshToken
}

// ResponseProblemHTTPCode 响应http错误码
func ResponseProblemHTTPCode(ctx iris.Context, code int, err error) {
	ctx.Application().Logger().Error(ctx.Path(), ": ", err)
	ctx.StatusCode(code)
	ctx.StopExecution()
	return
}

// BuildCondition 使用map构造GORM的where条件
//
// 用法：
//
// conditionString, values, _ := app.BuildCondition(map[string]interface{}{
// 	"itemNO":        "02WZ05133",
// 	"itemName like": "%22220",
// 	"id in":         []int{20, 19, 30},
// 	"num !=" : 20,
// })
//
// var student model.Student
//
// db.where(conditionString, values...).Find(&student)
func BuildCondition(where map[string]interface{}) (whereSQL string, values []interface{}, err error) {
	for key, value := range where {
		conditionKey := strings.Split(key, " ")
		if len(conditionKey) > 2 {
			return "", nil, fmt.Errorf("" + "map构建的条件格式不对，类似于'age >'")
		}
		if whereSQL != "" {
			whereSQL += " AND "
		}
		switch len(conditionKey) {
		case 1:
			whereSQL += fmt.Sprint(conditionKey[0], " = ?")
			values = append(values, value)
			break
		case 2:
			field := conditionKey[0]
			switch conditionKey[1] {
			case "=":
				whereSQL += fmt.Sprint(field, " = ?")
				values = append(values, value)
				break
			case ">":
				whereSQL += fmt.Sprint(field, " > ?")
				values = append(values, value)
				break
			case ">=":
				whereSQL += fmt.Sprint(field, " >= ?")
				values = append(values, value)
				break
			case "<":
				whereSQL += fmt.Sprint(field, " < ?")
				values = append(values, value)
				break
			case "<=":
				whereSQL += fmt.Sprint(field, " <= ?")
				values = append(values, value)
				break
			case "in":
				whereSQL += fmt.Sprint(field, " in (?)")
				values = append(values, value)
				break
			case "like":
				whereSQL += fmt.Sprint(field, " like ?")
				values = append(values, value)
				break
			case "<>":
				whereSQL += fmt.Sprint(field, " != ?")
				values = append(values, value)
				break
			case "!=":
				whereSQL += fmt.Sprint(field, " != ?")
				values = append(values, value)
				break
			}
			break
		}
	}
	return
}
