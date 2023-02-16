package es

import "testing"

func TestCreateIndex(t *testing.T) {
	var mapping = `{
		"mappings": {
		  "properties": {
			"age": {
			  "type": "integer"
			},
			"email": {
			  "type": "keyword"
			},
			"name": {
			  "type": "text"
			}
		  }
		}
	  }`
	if err := CreateIndex("index-test", mapping); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("create index-test ok")
	}
}

func TestCreateUpdateDocument(t *testing.T) {
	var document = `{
		"age": 30,
		"email":"36666@qq.com",
		"name": "小明3"
	}`
	if err := CreateUpdateDocument("index-test", "", document); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("create document ok")
	}
}

func TestCreateDocument(t *testing.T) {
	var document = `{
		"age": 30,
		"email":"36666@qq.com",
		"name": "小明2"
	}`
	if err := CreateDocument("index-test", "2", document); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("create document ok")
	}
}

func TestGetDocument(t *testing.T) {
	if res, err := GetDocument("index-test", "2"); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Get Document: ", res)
	}
}

func TestUpdateDocument(t *testing.T) {
	var document = `{
		"doc": {
			"age": 31,
			"email":"4444@qq.com",
			"name": "小明aa"
		}
	}`
	if err := UpdateDocument("index-test", "2", document); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("update ok")
	}
}

func TestUpdateDocumentsByQuery(t *testing.T) {
	var queryBody = `{
		"script": {
			"source": "ctx._source.name=params.name",
			"params": {
				"name":"newName"
			}
		},
		"query": {
			"term": {
				"email": "fff@qq.com"
			}
		}
	}`
	if err := UpdateDocumentsByQuery("index-test", queryBody); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Update Documents ok")
	}
}

func TestDeleteDocument(t *testing.T) {
	if err := DeleteDocument("index-test", "sC40T4YBcjhCCtcaHxK1"); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Delete ok")
	}
}

func TestDeleteDocumentsByQuery(t *testing.T) {
	// var queryBody = `{
	// 	"query":{
	// 		"match_all":{}
	// 	}
	// }`
	var queryBody = `{
		"query":{
			"match":{
				"age": 30
			}
		}
	}`
	if err := DeleteDocumentsByQuery("index-test", queryBody); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Delete Documents ok")
	}
}

func TestSearchDocuments(t *testing.T) {
	var queryBody = `{
		"query":{
			"match_all":{}
		}
	}`
	// var queryBody = `{
	// 	"query":{
	// 		"match":{
	// 			"name":"小明"
	// 		}
	// 	}
	// }`
	if list, total, err := SearchDocuments("index-test", queryBody, []string{"age:desc"}, 0, 2); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Search Documents: ", total, list)
	}
}
