/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:06:33
 * @LastEditTime: 2023-05-06 14:35:44
 */
package global

import (
	"context"
	"fmt"
	"iris-project/app/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDb     *mongo.Database
)

func init() {
	if !config.MongoDb.On {
		return
	}
	var (
		err error
		uri = fmt.Sprintf("mongodb://%v:%v@%v:%v/?tz=Asia/Shanghai", config.MongoDb.Username, config.MongoDb.Password, config.MongoDb.Host, config.MongoDb.Port) // "mongodb://admin:123@localhost:27017"
	)
	// log.Println(uri)
	uri = "mongodb+srv://356126067:hyl123456@cluster0.ru39nsc.mongodb.net/test"
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err.Error())
		for {
			time.Sleep(60 * time.Second) // 等待60秒再次连接
			log.Println("再次连接mongodb")
			MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
			if err == nil {
				break
			} else {
				log.Println(err.Error())
			}
		}
	}

	MongoDb = MongoClient.Database(config.MongoDb.Db)
}
