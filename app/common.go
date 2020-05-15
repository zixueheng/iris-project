package app

import (
	"iris-project/global"

	"github.com/go-playground/validator/v10"
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

	err := global.Validate.Struct(obj)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(global.ValidateTrans) {
			if len(e) > 0 {
				// p.Ctx.StatusCode(iris.StatusOK)
				// _, _ = ctx.JSON(APIData(false, nil, e))
				// return

				errmsg = e
				break
			}
		}
	}
	return
}
