/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2024-07-30 10:03:09
 */
package controller

import (
	"iris-project/app"
	adminmodel "iris-project/app/adminapi/model"
	"iris-project/app/dao"
	appmodel "iris-project/app/model"
	"iris-project/app/service"
	"iris-project/lib/util"
	"path"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
)

// File 控制器
type File struct {
	Ctx           iris.Context
	AuthAdminUser *adminmodel.AdminUser
}

// BeforeActivation 前置方法
func (c *File) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// @Tags        文件
// @Summary		单文件上传
// @Description	单文件上传
// @Accept		multipart/form-data
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		uploadfile	formData	file			true	"文件"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/upload [post]
func (c *File) PostFileUpload() {
	// Get the file from the request.
	file, info, err := c.Ctx.FormFile("uploadfile")
	if err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	defer file.Close()

	categoryID := c.Ctx.FormValueDefault("category_id", "0")
	onlyURL := c.Ctx.FormValueDefault("only_url", "1")

	f, err := service.Upload(file, info, uint32(util.ParseInt(categoryID)), "admin", c.AuthAdminUser.ID, "")

	if err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	// fname := util.GenFileName(path.Ext(info.Filename))
	// savePath := "/upload/" + util.TimeFormat(time.Now(), config.App.Dateformat) + "/"
	// fileutil.CreateFile("." + savePath)
	// fullPath := savePath + fname
	// if _, err = c.Ctx.SaveFormFile(info, "."+fullPath); err != nil {
	// 	c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
	// 	return
	// }

	if onlyURL == "1" {
		// URL ...
		type URL struct {
			URL string `json:"url"`
		}

		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", &URL{URL: f.URL}))
	} else {
		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", f))
	}

}

// @Tags        文件
// @Summary		多文件上传
// @Description	多文件上传
// @Accept		multipart/form-data
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		uploadfiles	formData	file			true	"文件"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/mutiple/upload [post]
func (c *File) PostFileMutipleUpload() {
	// Get the max post value size passed via iris.WithPostMaxMemory.
	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err := c.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		c.Ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	form := c.Ctx.Request().MultipartForm
	if form == nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "没有文件上传", nil))
		return
	}
	categoryID := c.Ctx.FormValueDefault("category_id", "0")
	onlyURL := c.Ctx.FormValueDefault("only_url", "1")

	// for k, v := range form.File {
	// 	fmt.Println(k, v)
	// }
	infos, ok := form.File["uploadfiles"]
	if !ok {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "没有文件上传", nil))
		return
	}
	fs, err := service.UploadMore(infos, uint32(util.ParseInt(categoryID)), "admin", c.AuthAdminUser.ID)
	if err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	if onlyURL == "1" {
		var urls = make([]string, 0)
		for _, f := range fs {
			urls = append(urls, f.URL)
		}

		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", urls))
	} else {
		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", fs))
	}
}

// @Tags        文件
// @Summary		文件列表
// @Description	文件列表
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		category_id     query       string  false	"分类ID"
// @Param		type            query       string  false	"类型"
// @Param		name            query       string  false	"名称"
// @Param		is_favor        query       int     false	"是否收藏"
// @Param		page            query       string  false	"页码"
// @Param		size            query       int     false	"页大小"
// @Success		200		{object}	appmodel.File		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/list [get]
func (c *File) GetFileList() {
	page, size := app.Pager(c.Ctx)

	var where = make(map[string]interface{})
	if categoryID := c.Ctx.URLParamIntDefault("category_id", -1); categoryID != -1 {
		// if categoryID == 0 {
		// 	where["category_id"] = 0
		// } else {
		// 	var categoryIDs []uint32
		// 	categoryIDs = append(categoryIDs, uint32(categoryID))
		// 	appmodel.GetFileCategoryChildIDs(&categoryIDs, uint32(categoryID))
		// 	where["category_id in"] = categoryIDs
		// }
		where["category_id"] = categoryID
	}
	if Type := c.Ctx.URLParamDefault("type", ""); Type != "" {
		where["type"] = Type
	}
	if name := c.Ctx.URLParamDefault("name", ""); name != "" {
		where["name like"] = "%" + name + "%"
	}
	if isFavor := c.Ctx.URLParamIntDefault("is_favor", 0); isFavor == 1 {
		var fileIDs = make([]uint32, 0)
		dao.Pluck(nil, &appmodel.FileFavor{}, map[string]interface{}{"admin_user_id": c.AuthAdminUser.ID}, "file_id", &fileIDs)
		if len(fileIDs) > 0 {
			where["id in"] = fileIDs
		} else {
			where["id"] = -1
		}
	}
	where["upload_by"] = "admin"

	var (
		list  []appmodel.File
		total int64
	)

	searchList := &dao.SearchListData{
		Where:    where,
		Fields:   []string{},
		OrderBys: []string{"id desc"},
		Preloads: []string{},
		Page:     page,
		Size:     size,
	}
	if err := searchList.GetList(&list, &total); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	for i, v := range list {
		v.LoadRelatedField(c.AuthAdminUser.ID)
		list[i] = v
	}

	var (
		statistic = make(map[string]interface{})
		totalSize int64
	)
	dao.Pluck(nil, &appmodel.File{}, nil, "COALESCE(SUM(size), 0) as total_size", &totalSize)
	statistic["total_size"] = totalSize

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: list, Total: total, Statistic: statistic}))
}

// @Tags        文件
// @Summary		收藏或取消收藏文件
// @Description	收藏或取消收藏文件
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		request	body		string	true	"请求体，格式：{file_id:1}"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/favor [put]
func (c *File) PutFileFavor() {
	type Request struct {
		FileID uint32 `json:"file_id" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	var favor = &appmodel.FileFavor{}
	dao.FindOne(nil, favor, map[string]interface{}{"file_id": param.FileID, "admin_user_id": c.AuthAdminUser.ID})
	if favor.FileID == 0 {
		favor.FileID = param.FileID
		favor.AdminUserID = c.AuthAdminUser.ID
		dao.SaveOne(nil, favor)
	} else {
		dao.DeleteAll(nil, favor, map[string]interface{}{"file_id": param.FileID, "admin_user_id": c.AuthAdminUser.ID})
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        文件
// @Summary		修改文件名
// @Description	修改文件名
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		request	body		string	true	"请求体，格式：{file_id:1,name:新名称}"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/name [put]
func (c *File) PutFileName() {
	type Request struct {
		FileID uint32 `json:"file_id" validate:"required"`
		Name   string `json:"name" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	var file = &appmodel.File{ID: param.FileID}
	if !dao.GetByID(file) {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "文件不存在", nil))
		return
	}

	if err := dao.UpdateByID(&appmodel.File{ID: param.FileID}, map[string]interface{}{"name": param.Name + path.Ext(file.Name)}); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        文件
// @Summary		更改文件分类
// @Description	更改文件分类
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		request	body		string	true	"请求体，格式：{ids:[1,2],name:新名称}"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file/category [put]
func (c *File) PutFileCategory() {
	type Request struct {
		IDs        []uint32 `json:"ids" validate:"required,min=1"`
		CategoryID uint32   `json:"category_id" validate:"-"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	if err := dao.UpdateAll(nil, &appmodel.File{}, map[string]interface{}{"id in": param.IDs}, map[string]interface{}{"category_id": param.CategoryID}); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        文件
// @Summary		删除文件
// @Description	删除文件
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		request	body		string	true	"请求体，格式：{ids:[1,2]}"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/file [delete]
func (c *File) DeleteFile() {
	type Request struct {
		IDs []uint32 `json:"ids" validate:"required,min=1"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	deleted := service.DeleteFiles(param.IDs)
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", map[string]interface{}{"deleted": deleted}))
}

// @Tags        文件
// @Summary		文件分类树
// @Description	文件分类树
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Success		200		{object}	appmodel.FileCategoryTree		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/filecategory/tree [get]
func (c *File) GetFilecategoryTree() {
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", appmodel.GetFileCategoryTree()))
}

// @Tags        文件
// @Summary		文件分类选择项
// @Description	文件分类选择项
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Success		200		{object}	appmodel.FileCategory		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/filecategory/select [get]
func (c *File) GetFilecategorySelect() {
	var list []*appmodel.FileCategory
	dao.FindAll(nil, &list, nil)

	var box = make([]*appmodel.FileCategory, 0)
	var root = &appmodel.FileCategory{ID: 0, Name: "顶级按钮", Level: 0, HTML: ""}
	box = append(box, root)
	var fn func([]*appmodel.FileCategory, *appmodel.FileCategory)
	fn = func(ls []*appmodel.FileCategory, node *appmodel.FileCategory) {
		for _, v := range ls {
			if v.PID == node.ID {
				v.Level = node.Level + 1
				v.HTML = strings.Repeat("|----", v.Level)
				box = append(box, v)
				fn(ls, v)
			}
		}
	}
	fn(list, root)

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", box))
}

// @Tags        文件
// @Summary		创建或更新文件分类
// @Description	创建或更新文件分类
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		account	body		appmodel.FileCategory	true	"请求体"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/filecategory [post]
func (c *File) PostFilecategory() {
	var category dao.Model = &appmodel.FileCategory{}

	errmsg := app.CheckRequest(c.Ctx, category)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	if err := dao.CreateUpdate(category); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        文件
// @Summary		删除文件分类
// @Description	删除文件分类
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"token"
// @Param		id path int true	"ID"
// @Success		200		{object}	model.AdminUser		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/filecategory/{id} [delete]
func (c *File) DeleteFilecategoryBy(id uint32) {
	var categoryIDs []uint32
	categoryIDs = append(categoryIDs, id)
	appmodel.GetFileCategoryChildIDs(&categoryIDs, id)

	if err := dao.Transaction(func(tx *gorm.DB) error {
		if err := dao.UpdateAll(tx, &appmodel.File{}, map[string]interface{}{"category_id in": categoryIDs}, map[string]interface{}{"category_id": 0}); err != nil {
			return err
		}
		if err := dao.DeleteAll(tx, &appmodel.FileCategory{}, map[string]interface{}{"id in": categoryIDs}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}
