/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:04:43
 * @LastEditTime: 2023-04-20 17:26:14
 */
package mongodb

import (
	"errors"
	"iris-project/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	GetID() primitive.ObjectID
	SetID(primitive.ObjectID)
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

// GetByID 根据ID查询
func GetByID(db *mongo.Database, m Model) bool {
	if m.GetID().IsZero() {
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

// QueryOpts 查询选项
type QueryOpts struct {
	Projection bson.D // 包含或排除字段，1包含 0排除，如 bson.D{{"course_id", 0}, {"enrollment", 0}}
	Sort       bson.D // 排序，可按多个字段排序 1 asc -1 desc，如 bson.D{{"total_price", 1}
	Limit      int64  // 限定数量
}

// FindOne 查找一个结果
//
// filter 查询条件：bson.D{{"uid", 1}}
func FindOne(db *mongo.Database, m Model, filter interface{}, op *QueryOpts) error {
	if db == nil {
		db = GetDB()
	}

	var opts = make([]*options.FindOneOptions, 0)

	if op != nil {
		if op.Projection != nil {
			opts = append(opts, options.FindOne().SetProjection(op.Projection))
		}
		if op.Sort != nil {
			opts = append(opts, options.FindOne().SetSort(op.Sort))
		}
	}

	if err := db.Collection(m.GetName()).FindOne(nil, filter, opts...).Decode(m); err == mongo.ErrNoDocuments {
		return err
	}

	return nil
}

// FindAll 查找一个结果
//
// filter 查询条件：bson.D{{"uid", 1}}
func FindAll(db *mongo.Database, results []Model, filter interface{}, op *QueryOpts) error {
	if db == nil {
		db = GetDB()
	}

	var name = ""
	if len(results)==1{
		name = results[0].GetName()
		results = results[:0]
	}else{
		return errors.New("results should have one element")
	}
	

	var opts = make([]*options.FindOptions, 0)

	if op != nil {
		if op.Projection != nil {
			opts = append(opts, options.Find().SetProjection(op.Projection))
		}
		if op.Sort != nil {
			opts = append(opts, options.Find().SetSort(op.Sort))
		}
		if op.Limit != 0 {
			opts = append(opts, options.Find().SetLimit(op.Limit))
		}
	}

	if cursor, err := db.Collection(name).Find(nil, filter, opts...); err == nil {
		return err
	} else {
		if err = cursor.All(nil, &results); err != nil {
			return err
		}
	}

	return nil
}

func CreateUpdate(db *mongo.Database, m Model) error {
	if db == nil {
		db = GetDB()
	}
	if m.GetID().IsZero() { // create
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

// SaveOne 保存一个
func SaveOne(db *mongo.Database, m Model) error {
	if db == nil {
		db = GetDB()
	}

	if res, err := db.Collection(m.GetName()).InsertOne(nil, m); err != nil {
		return err
	} else {
		// log.Printf("结果：%+v\n", res)
		id, _ := res.InsertedID.(primitive.ObjectID)
		m.SetID(id)
	}
	return nil
}

// SaveAll 保存多个
func SaveAll(db *mongo.Database, slice []interface{}) error {
	if len(slice) == 0 {
		return errors.New("save empty slice")
	}

	if db == nil {
		db = GetDB()
	}

	var name = ""
	for _, v := range slice {
		if m, ok := v.(Model); !ok {
			return errors.New("type assert failed")
		} else {
			name = m.GetName()
		}
	}

	if res, err := db.Collection(name).InsertMany(nil, slice); err != nil {
		return err
	} else {
		for i, insertedID := range res.InsertedIDs {
			id, _ := insertedID.(primitive.ObjectID)
			if t, ok := slice[i].(Model); ok {
				t.SetID(id)
			}
		}
	}
	return nil
}
