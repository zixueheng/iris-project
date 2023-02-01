/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-05-08 15:21:44
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/config"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Version 控制器
type Version struct {
	Ctx iris.Context
}

// AfterActivation 后置方法
func (c *Version) AfterActivation(a mvc.AfterActivation) {

}

// PostVersionCompare 最新版本
func (c *Version) PostVersionCompare() {
	var param = struct {
		Version  string `json:"version"`
		Platform string `json:"platform"`
	}{}

	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}

	if param.Version == "" {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "请输入当前APP版本号", nil))
		return
	}
	var (
		v1 string
		v2 string
	)

	v1 = strings.TrimSpace(param.Version)
	v1Slice := strings.Split(v1, ".")
	if len(v1Slice) != 3 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "APP版本号格式不正确，例：1.0.0", nil))
		return
	}

	version := config.GetVersion()
	// fmt.Println(*version)
	if param.Platform == "ios" {
		v2 = strings.TrimSpace(version.Ios.Version)
	} else {
		v2 = strings.TrimSpace(version.Android.Version)
	}

	v2Slice := strings.Split(v2, ".")
	if len(v2Slice) != 3 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "系统设置的最新版本号不正确，请联系管理员", nil))
		return
	}
	type VersionRes struct {
		LatestVersion string   `json:"latest_version"`
		Title         string   `json:"title"`
		Msgs          []string `json:"msgs"`
		Downloadurl   string   `json:"download_url"`
	}
	var res *VersionRes

	for index := range v1Slice {
		if v1Slice[index] < v2Slice[index] {
			if param.Platform == "ios" {
				res = &VersionRes{
					LatestVersion: v2,
					Title:         version.Ios.Title,
					Msgs:          version.Ios.Msgs,
					Downloadurl:   version.Ios.Downloadurl,
				}
			} else {
				res = &VersionRes{
					LatestVersion: v2,
					Title:         version.Android.Title,
					Msgs:          version.Android.Msgs,
					Downloadurl:   version.Android.Downloadurl,
				}
			}
			break
		}
	}

	if res == nil {
		if param.Platform == "ios" {
			res = &VersionRes{
				LatestVersion: version.Ios.Version,
				Title:         "您当前是最新版本",
			}
		} else {
			res = &VersionRes{
				LatestVersion: version.Android.Version,
				Title:         "您当前是最新版本",
			}
		}
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", res))

}
