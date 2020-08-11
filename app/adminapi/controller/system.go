package controller

import (
	"bufio"
	"io"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/config"
	"iris-project/lib/util"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// System 控制器
type System struct {
	Ctx           iris.Context     // IRIS框架会自动注入 Context
	AuthAdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (s *System) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

// GetSystemDbBackup 数据库备份
func (s *System) GetSystemDbBackup() {
	pathname, err := util.BackupMySQLDb(config.DB.Host, util.ParseString(int(config.DB.Port)), config.DB.User, config.DB.Password, config.DB.Name, "", "./upload/backup/")
	if err != nil {
		s.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
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
	s.Ctx.Header("Content-Disposition", "attachment; filename="+path.Base(pathname))
	s.Ctx.Header("Content-Type", "text/plain")
	// s.Ctx.Write(data)

	s.Ctx.Header("Transfer-Encoding", "chunked")

	s.Ctx.StreamWriter(func(w io.Writer) bool {
		file, _ := os.Open(pathname)
		if err != nil {
			return false
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 读到一个换行就结束
			if err == io.EOF {                  // io.EOF 表示文件的末尾
				break
			}
			// fmt.Print(str)
			s.Ctx.WriteString(str)
		}
		return false
	})
}
