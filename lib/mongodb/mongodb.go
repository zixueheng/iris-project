/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:04:43
 * @LastEditTime: 2023-03-15 17:39:28
 */
package mongodb

import (
	"errors"
	"iris-project/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	GetID() string
	GetName() string
}

// 获取Mongodb客户端
func GetClient() *mongo.Client {
	return global.MongoClient
}

// 获取默认数据库
func GetDB() *mongo.Database {
	return global.MongoDb
}

// 切换其他库
func SwitchDB(name string) *mongo.Database {
	return GetClient().Database(name)
}

func GetByID(db *mongo.Database, m Model) bool {
	if m.GetID() == "" {
		return false
	}

	if db == nil {
		db = GetDB()
	}

	if err := db.Collection(m.GetName()).FindOne(nil, bson.D{{"_id", m.GetID()}}).Decode(m); err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func CreateUpdate(db *mongo.Database, m Model) error {
	if db == nil {
		db = GetDB()
	}
	if m.GetID() == "" { // create
		if _, err := db.Collection(m.GetName()).InsertOne(nil, m); err != nil {
			return err
		}
	} else { // update
		if _, err := db.Collection(m.GetName()).ReplaceOne(nil, bson.D{{"_id", m.GetID()}}, m); err != nil {
			return err
		}
	}
	return nil
}

func SaveAll(db *mongo.Database, slice []interface{}) error {
	if len(slice) == 0 {
		return errors.New("Save empty slice")
	}

	if db == nil {
		db = GetDB()
	}

	m, ok := slice[0].(Model)
	if !ok {
		return errors.New("Type assert failed")
	}

	if _, err := db.Collection(m.GetName()).InsertMany(nil, slice); err != nil {
		return err
	}
	return nil
}
