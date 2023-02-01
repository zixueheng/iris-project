/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-11 16:41:22
 * @LastEditTime: 2021-08-28 18:16:24
 */
package service

import (
	"errors"
	"io"
	"iris-project/app/config"
	"iris-project/app/dao"
	appmodel "iris-project/app/model"
	fileutil "iris-project/lib/file"
	"iris-project/lib/util"
	"mime/multipart"
	"os"
	"path"
	"time"
)

// Upload 上传文件
func Upload(file multipart.File, info *multipart.FileHeader, categoryID uint32, uploadBy string, uploadByID uint32, newfilename string) (*appmodel.File, error) {
	if info.Size > config.Upload.Maxsize<<20 {
		return nil, errors.New("上传文件不能超过" + util.ParseString(int(config.Upload.Maxsize)) + "M")
	}
	var (
		domain, fullpath, oldname, filename string
		err                                 error
	)
	if newfilename == "" {
		oldname = info.Filename
	} else {
		oldname = newfilename
	}

	switch config.Upload.Storage {
	case "local":
		domain, fullpath, filename, err = uploadLocal(file, info)
	case "qiniu":
		domain, fullpath, filename, err = uploadQiniu(file, info)
	case "tencent":
		domain, fullpath, filename, err = uploadTencent(file, info)
	case "aliyun":
		domain, fullpath, filename, err = uploadAliyun(file, info)
	}
	if err != nil {
		return nil, err
	}
	fileObj := &appmodel.File{
		Name:       oldname,
		FileName:   filename,
		Path:       fullpath,
		URL:        domain + fullpath,
		Size:       info.Size,
		Type:       appmodel.FileType(path.Ext(info.Filename)),
		FileMime:   info.Header.Get("Content-Type"),
		CategoryID: categoryID,
		Storage:    config.Upload.Storage,
		UploadBy:   uploadBy,
		UploadByID: uploadByID,
	}
	dao.SaveOne(nil, fileObj)
	return fileObj, nil
}

// UploadMore 上传多个文件
func UploadMore(infos []*multipart.FileHeader, categoryID uint32, uploadBy string, uploadByID uint32) (files []*appmodel.File, err error) {
	for _, info := range infos {
		if info.Size > config.Upload.Maxsize<<20 {
			return nil, errors.New("上传文件不能超过" + util.ParseString(int(config.Upload.Maxsize)) + "M")
		}
	}
	for _, info := range infos {
		file, _ := info.Open()
		f, _ := Upload(file, info, categoryID, uploadBy, uploadByID, "")
		files = append(files, f)
		file.Close()
	}
	return
}

// uploadLocal 本地上传
func uploadLocal(file multipart.File, info *multipart.FileHeader) (domain, fullpath, filename string, err error) {
	domain = config.Upload.Local.Domain
	filename = util.GenFileName(path.Ext(info.Filename))
	savePath := "/upload/" + util.TimeFormat(time.Now(), config.App.Dateformat) + "/"

	fileutil.CreateFile("." + savePath)
	fullpath = savePath + filename
	out, err := os.OpenFile("."+fullpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return
	}
	return
}

// uploadQiniu 七牛云上传
func uploadQiniu(file multipart.File, info *multipart.FileHeader) (domain, fullpath, filename string, err error) {
	return
}

// uploadTencent 腾讯云上传
func uploadTencent(file multipart.File, info *multipart.FileHeader) (domain, fullpath, filename string, err error) {
	return
}

// uploadAliyun 阿里云上传
func uploadAliyun(file multipart.File, info *multipart.FileHeader) (domain, fullpath, filename string, err error) {
	return
}

// DeleteFiles 删除多个文件
func DeleteFiles(ids []uint32) (deleted uint32) {
	var files []*appmodel.File
	dao.FindAll(nil, &files, map[string]interface{}{"id in": ids})
	for _, file := range files {
		if err := DeleteFile(file); err == nil {
			deleted++
		}
	}
	return
}

// DeleteFile 删除文件
func DeleteFile(file *appmodel.File) (err error) {
	switch file.Storage {
	case "local":
		err = deleteLocalFile(file.Path)
	case "qiniu":
		err = deleteQiniuFile(file.Path)
	case "tencent":
		err = deleteTencentFile(file.Path)
	case "aliyun":
		err = deleteAliyunFile(file.Path)
	}
	// if err == nil {
	// 	dao.DeleteByID(file)
	// }
	dao.DeleteByID(file)
	return
}

// deleteLocalFile 删除本地文件
func deleteLocalFile(path string) (err error) {
	err = os.Remove("." + path)
	return
}

// deleteQiniuFile 删除七牛云文件
func deleteQiniuFile(path string) (err error) {
	return
}

// deleteTencentFile 删除腾讯云文件
func deleteTencentFile(path string) (err error) {
	return
}

// deleteAliyunFile 删除阿里云文件
func deleteAliyunFile(path string) (err error) {
	return
}
