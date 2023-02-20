/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-11 10:30:51
 * @LastEditTime: 2023-02-20 15:34:25
 */
package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"iris-project/global"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// GetEsClient 获取ES连接
func GetEsClient() *elasticsearch.Client {
	return global.EsClient
}

// 检查响应，返回错误，如果后续不需要使用响应则应该关闭它
func checkResponse(res *esapi.Response, closeBody bool) error {
	if closeBody {
		defer res.Body.Close()
	}
	if res.IsError() {
		log.Println("Elasticsearch: ", res.Status())
		// log.Printf("RES: %+v", res)
		// var bytes []byte
		// if _, err := res.Body.Read(bytes); err != nil {
		// 	log.Println(err.Error())
		// }
		return errors.New(res.String())
	}
	return nil
}

// 返回Body内容
func getResponse(res *esapi.Response) string {
	var (
		out = new(bytes.Buffer)
	)

	if res != nil && res.Body != nil {
		defer res.Body.Close()
		if _, err := out.ReadFrom(res.Body); err != nil {
			out.WriteString("")
			return out.String()
		}
	}

	return out.String()
}

// IndexExist 检查Index是否存在，true存在、false不存在
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-exists.html
func IndexExist(index string) bool {
	var (
		res *esapi.Response
		err error
	)
	if res, err = GetEsClient().Indices.Exists(
		[]string{index},
	); err != nil {
		return false
	}
	if checkResponse(res, true) == nil {
		return true
	}

	return false
}

// CreateIndex 创建索引
// index 索引名
// mapping 映射
//
//	映射定义：https://www.elastic.co/guide/en/elasticsearch/reference/master/explicit-mapping.html
//	字段类型：https://www.elastic.co/guide/en/elasticsearch/reference/master/mapping-types.html
func CreateIndex(index, mapping string) (err error) {
	var res *esapi.Response
	if res, err = GetEsClient().Indices.Create(
		index,
		GetEsClient().Indices.Create.WithBody(strings.NewReader(mapping)),
	); err != nil {
		return
	}

	return checkResponse(res, true)
}

// 创建修改文档（POST/PUT /<index>/_doc/<_id>）
//
// id 已存在会修改文档
//
// id 不传会生成随机ID
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func CreateUpdateDocument(index string, id, document string) (err error) {
	var res *esapi.Response
	if res, err = GetEsClient().Index(index,
		strings.NewReader(document),
		GetEsClient().Index.WithDocumentID(id),
		GetEsClient().Index.WithRefresh("true"),
	); err != nil {
		return
	}
	return checkResponse(res, true)
}

// 创建文档（PUT /<index>/_create/id）
//
// id 必传，id已存在会报错
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func CreateDocument(index string, id, document string) (err error) {
	if id == "" {
		return errors.New("ID必传")
	}
	var res *esapi.Response
	if res, err = GetEsClient().Create(index,
		id,
		strings.NewReader(document),
		GetEsClient().Create.WithRefresh("true"),
	); err != nil {
		return
	}

	return checkResponse(res, true)
}

// 获取文档（GET /<index>/_source/<_id>）
//
// id 必传，id已存在会报错
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-get.html
func GetDocument(index string, id string) (document string, err error) {
	if id == "" {
		return "", errors.New("ID必传")
	}
	var res *esapi.Response
	if res, err = GetEsClient().GetSource(index,
		id,
	); err != nil {
		return
	}

	if err = checkResponse(res, false); err != nil {
		return
	}

	return getResponse(res), nil
}

// 修改文档（PUT /<index>/_update/id）
//
// id 必传，id已存在会报错；通常用来修改全部或部分字段，用法较多请查看文档：
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-update.html
func UpdateDocument(index string, id, document string) (err error) {
	if id == "" {
		return errors.New("ID必传")
	}
	var res *esapi.Response
	if res, err = GetEsClient().Update(index,
		id,
		strings.NewReader(document),
		GetEsClient().Update.WithRefresh("true"),
	); err != nil {
		return
	}
	return checkResponse(res, true)
}

// 指定条件更新文档（POST /<index>/_update_by_query）
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-update-by-query.html
func UpdateDocumentsByQuery(index string, queryBody string) (err error) {
	var res *esapi.Response
	if res, err = GetEsClient().UpdateByQuery(
		[]string{index},
		GetEsClient().UpdateByQuery.WithBody(strings.NewReader(queryBody)),
		GetEsClient().UpdateByQuery.WithRefresh(true),
	); err != nil {
		return
	}

	return checkResponse(res, true)
}

// 删除文档（DELETE /<index>/_doc/<_id>）
//
// id 必传，id已存在会报错
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-delete.html
func DeleteDocument(index string, id string) (err error) {
	if id == "" {
		return errors.New("ID必传")
	}
	var res *esapi.Response
	if res, err = GetEsClient().Delete(index,
		id,
		GetEsClient().Delete.WithRefresh("true"),
	); err != nil {
		return
	}

	return checkResponse(res, true)
}

// 指定条件删除文档（POST /<index>/_delete_by_query）
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-delete-by-query.html
func DeleteDocumentsByQuery(index string, queryBody string) (err error) {
	var res *esapi.Response
	if res, err = GetEsClient().DeleteByQuery(
		[]string{index},
		strings.NewReader(queryBody),
		GetEsClient().DeleteByQuery.WithRefresh(true),
	); err != nil {
		return
	}

	return checkResponse(res, true)
}

// 搜索文档（POST /<index>/_search/）
//
// sorts 排序 <field>:<direction>，
// from 从0开始，size最多10000
//
// https://www.elastic.co/guide/en/elasticsearch/reference/master/search-search.html
func SearchDocuments(index string, queryBody string, sorts []string, from, size int) (docments []*ResponseRecord, total int64, err error) {
	var res *esapi.Response
	if res, err = GetEsClient().Search(
		GetEsClient().Search.WithIndex(index),
		// GetEsClient().Search.WithSource(),
		GetEsClient().Search.WithBody(strings.NewReader(queryBody)),
		GetEsClient().Search.WithSort(sorts...),
		GetEsClient().Search.WithFrom(from),
		GetEsClient().Search.WithSize(size),
	); err != nil {
		return
	}

	if err = checkResponse(res, false); err != nil {
		return
	}

	var responseSearch = &ResponseSearch{}
	if err = json.Unmarshal([]byte(getResponse(res)), responseSearch); err != nil {
		return
	}

	return responseSearch.Hits.Hits, responseSearch.Hits.Total.Value, nil

}
