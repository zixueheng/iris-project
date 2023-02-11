package es

import (
	"iris-project/global"

	"github.com/elastic/go-elasticsearch/v8"
)

// GetEsClient 获取ES连接
func GetEsClient() *elasticsearch.Client {
	return global.EsClient
}

func CreateIndex(index string) {
}
