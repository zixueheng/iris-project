/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-04-28 14:59:03
 */
package controller

import (
	"bufio"
	"io"
	"iris-project/app"
	adminmodel "iris-project/app/adminapi/model"
	"iris-project/app/config"
	"iris-project/lib/util"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// System 控制器
type System struct {
	Ctx           iris.Context
	AuthAdminUser *adminmodel.AdminUser
}

// BeforeActivation 前置方法
func (c *System) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetSystemDbBackup 数据库备份
func (c *System) GetSystemDbBackup() {
	pathname, err := util.BackupMySQLDb(config.DB.Host, util.ParseString(int(config.DB.Port)), config.DB.User, config.DB.Password, config.DB.Name, "", "./upload/backup/")
	if err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	// path = config.App.Fronturl + path[1:]
	// res := struct {
	// 	Path string `json:"path"`
	// }{
	// 	Path: path,
	// }
	// s.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", res))

	// data, _ := ioutil.ReadFile(pathname)

	// request := pr.Ctx.Request()
	c.Ctx.Header("Content-Disposition", "attachment; filename="+path.Base(pathname))
	c.Ctx.Header("Content-Type", "text/plain")
	// s.Ctx.Write(data)

	c.Ctx.Header("Transfer-Encoding", "chunked")

	c.Ctx.StreamWriter(func(w io.Writer) error {
		file, err := os.Open(pathname)
		if err != nil {
			return err
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 读到一个换行就结束
			if err == io.EOF {                  // io.EOF 表示文件的末尾
				break
			}
			// fmt.Print(str)
			c.Ctx.WriteString(str)
		}
		return err
	})
}
