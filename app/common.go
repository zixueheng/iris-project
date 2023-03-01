/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2023-03-01 10:58:49
 */
package app

import (
	"fmt"
	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/util"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// APP公共函数

// Pager 获取请求参数分页
func Pager(ctx iris.Context) (page, size int) {
	page = ctx.URLParamIntDefault("page", 1)
	size = ctx.URLParamIntDefault("size", 10)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 50
	}
	return
}

// CheckRequest 检查请求参数，返回错误提示
func CheckRequest(ctx iris.Context, obj interface{}) (errmsg string) {
	// if err := ctx.ReadJSON(obj); err != nil {
	// 	errmsg = err.Error()
	// 	return
	// }
	if errmsg = ReadJSONFromBody(ctx, obj); errmsg != "" {
		return
	}

	errmsg = ValidateStruct(obj)
	return
}

// ReadJSONFromBody 从请求Body中获取数据
func ReadJSONFromBody(ctx iris.Context, obj interface{}) (errmsg string) {
	if err := ctx.ReadJSON(obj); err != nil {
		errmsg = err.Error()
		return
	}
	return
}

// ValidateStruct 验证结构体数据
func ValidateStruct(obj interface{}) string {
	var errmsg string = ""
	err := global.Validate.Struct(obj)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(global.ValidateTrans) {
			// fmt.Println(e)
			if len(e) > 0 {
				errmsg += " " + e
			}
		}
	}
	return errmsg
}

// GenTokenAndRefreshToken 生成Token和RefreshToken，并返回Token和RefreshToken的过期时间
//
// key JWT要保存的键名，id 键值，
// tokenMinutes Token多少分钟后过期，
// refreshTokenMinutes 刷新Token多少分钟后过期
// client 端
func GenTokenAndRefreshToken(key string, id int, tokenMinutes, refreshTokenMinutes int, client string) (token, refreshToken string, tokenExpired, refreshTokenExpired time.Time) {
	var now = time.Now()
	tokenExpired = now.Add(time.Minute * time.Duration(tokenMinutes))

	// 获取一个 Token，参数一：签名方法、参数二：要保存的数据
	tokenObj := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		key:              util.ParseString(id),
		global.ClientKey: client,
		"exp":            tokenExpired.Unix(), //util.TimeFormat(tokenExpired, ""),
		"iat":            now.Unix(),          //util.TimeFormat(now, ""),
	})

	// Sign and get the complete encoded token as a string using the secret
	token, _ = tokenObj.SignedString([]byte(config.App.Jwtsecret))

	refreshToken = util.GetRandomString(64)
	refreshTokenExpired = now.Add(time.Minute * time.Duration(refreshTokenMinutes))

	return
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
//	conditionString, values, _ := app.BuildCondition(map[string]interface{}{
//		"itemNO":        "02WZ05133",
//		"itemName like": "%22220",
//		"id in":         []int{20, 19, 30},
//		"num !=" : 20,
//	})
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
		case 2:
			field := conditionKey[0]
			switch conditionKey[1] {
			case "=":
				whereSQL += fmt.Sprint(field, " = ?")
				values = append(values, value)
			case ">":
				whereSQL += fmt.Sprint(field, " > ?")
				values = append(values, value)
			case ">=":
				whereSQL += fmt.Sprint(field, " >= ?")
				values = append(values, value)
			case "<":
				whereSQL += fmt.Sprint(field, " < ?")
				values = append(values, value)
			case "<=":
				whereSQL += fmt.Sprint(field, " <= ?")
				values = append(values, value)
			case "in":
				whereSQL += fmt.Sprint(field, " in (?)")
				values = append(values, value)
			case "like":
				whereSQL += fmt.Sprint(field, " like ?")
				values = append(values, value)
			case "<>":
				whereSQL += fmt.Sprint(field, " != ?")
				values = append(values, value)
			case "!=":
				whereSQL += fmt.Sprint(field, " != ?")
				values = append(values, value)
			}
		}
	}
	return
}

// GetFormStatus 获取统一表单状态文本
func GetFormStatus(status int8) string {
	switch status {
	case 1:
		return "编辑中"
	case 2:
		return "审核中"
	case 3:
		return "审核通过"
	case 4:
		return "审核失败"
	}
	return ""
}

// ExportExcel 导出Excel
func ExportExcel(ctx iris.Context, fileName string, title []string, data [][]interface{}, saveLocal bool) {
	var (
		f     = excelize.NewFile()
		sheet = "Sheet1"
	)

	f.SetSheetRow(sheet, "A1", &title)
	i := 2
	for _, row := range data {
		f.SetSheetRow(sheet, "A"+util.ParseString(i), &row)
		i++
	}

	if saveLocal {
		rootpath, _ := os.Getwd()
		pathname := filepath.Join(rootpath, "upload/excel/"+fileName)

		if err := f.SaveAs(pathname); err != nil {
			fmt.Println(err)
		}
	}

	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/vnd.ms-excel")
	// pr.Ctx.Write(excel)

	w := ctx.ResponseWriter()

	if _, err := f.WriteTo(w); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}
