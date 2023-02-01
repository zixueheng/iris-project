/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2021-04-28 14:48:37
 */
package global

import (
	"fmt"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 表单验证
// https://github.com/go-playground/validator
var (
	uni           *ut.UniversalTranslator
	Validate      *validator.Validate
	ValidateTrans ut.Translator
)

func init() {
	zh2 := zh.New()
	uni = ut.New(zh2, zh2)
	ValidateTrans, _ = uni.GetTranslator("zh")
	Validate = validator.New()
	// 收集结构体中的comment标签，用于替换英文字段名称，这样返回错误就能展示中文字段名称了
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		comment := fld.Tag.Get("comment")
		if comment == "" {
			return fld.Name
		}
		return comment
		// return fld.Tag.Get("comment")
	})
	if err := zh_translations.RegisterDefaultTranslations(Validate, ValidateTrans); err != nil {
		fmt.Println(fmt.Sprintf("RegisterDefaultTranslations %v", err))
	}
}
