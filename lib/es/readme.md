<!--
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-11 15:33:22
 * @LastEditTime: 2023-02-11 17:11:33
-->
# ElasticSearch API
## 索引
创建索引shopping
http://127.0.0.1:9200/shopping PUT

查询索引
http://127.0.0.1:9200/shopping GET 

查看所有索引
http://127.0.0.1:9200/_cat/indices?v GET 

删除索引
http://127.0.0.1:9200/shopping DELTET

## 文档操作
文档创建
http://127.0.0.1:9200/shopping/_doc POST
```json
{
    "title":"小米手机",
    "category":"小米",
    "images":"http://www.gulixueyuan.com/xm.jpg",
    "price":3999.00
}
```

指定ID方式创建文档
http://127.0.0.1:9200/shopping/_doc/1 POST
```json
{
    "title":"小米手机",
    "category":"小米",
    "images":"http://www.gulixueyuan.com/xm.jpg",
    "price":3999.00
}
```

文档全量修改
http://127.0.0.1:9200/shopping/_doc/1 POST
```json
{
    "title":"小米手机2",
    "category":"小米",
    "images":"http://www.gulixueyuan.com/xm.jpg",
    "price":3999.00
}
```
文档局部修改（只修改部分字段）
http://127.0.0.1:9200/shopping/_update/1
```json
{
	"doc": {
		"title":"小米手机",
		"category":"小米"
	}
}
```

按ID删除文档
http://127.0.0.1:9200/shopping/_doc/1 DELTET

## 文档查询
指定ID查询
http://127.0.0.1:9200/shopping/_doc/1 GET 

查询所有
http://127.0.0.1:9200/shopping/_search GET

URL参数查询（不常使用）
http://127.0.0.1:9200/shopping/_search?q=category:小米 GET

请求体带参查询
http://127.0.0.1:9200/shopping/_search
```json
{
	"query":{
		"match":{
			"category":"小米" // 查询的字段
		}
	},

    // 查询所有
	"query":{
        "match_all":{} // 查询所有 注意和`match`二选一
	},

    // 多条件查询 
    "query":{
        // 多条件查询 
        "bool": { // 找出小米牌子并且价格为3999元的
            "must": [ // must相当于数据库的and
                {
                    "match": {
                        "category": "小米"
                    }
                },
                {
                    "match": {
                        "price": 3999.00
                    }
                }
            ]
        }
	},

     // 多条件查询 
    "query":{
        // 多条件查询 
        "bool": { // 找出小米牌子或者价格为3999元的
            "should": [ // should相当于数据库的or
                {
                    "match": {
                        "category": "小米"
                    }
                },
                {
                    "match": {
                        "price": 3999.00
                    }
                }
            ]
        }
	},

    // 找出小米和华为的牌子，价格大于2000（有问题）
    "query": {
        "bool": {
            "should": [ //  should 没有效果，
                {
                    "match": {
                        "category": "小米"
                    }
                },
                {
                    "match": {
                        "category": "华为"
                    }
                }
            ],
            "filter": {
                "range": {
                    "price": {
                        "gt": 2000
                    }
                }
            }
        }
    },


    
    "_source":["title"], // 返回结果指定字段
    "from": 0, // 起始下标
    "size": 10, // 查询长度
    "sort":{ // 排序
		"price":{
			"order":"desc"
		}
	}
}
```

## 映射关系
先创建一个索引
http://127.0.0.1:9200/user PUT

再创建映射(定义字段)
http://127.0.0.1:9200/user/_mapping PUT
```json
{
    "properties": {
        "name":{
        	"type": "text",
        	"index": true
        },
        "sex":{
        	"type": "keyword",
        	"index": true
        },
        "tel":{
        	"type": "keyword",
        	"index": false
        }
    }
}
```

查询映射
http://127.0.0.1:9200/user/_mapping GET

增加数据
http://127.0.0.1:9200/user/_create/1001 PUT
{
	"name":"小米",
	"sex":"男的",
	"tel":"1111"
}
