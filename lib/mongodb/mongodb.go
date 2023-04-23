/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:04:43
 * @LastEditTime: 2023-04-23 17:04:58
 */
package mongodb

import (
	"errors"
	"iris-project/global"
	"iris-project/lib/util"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	GetID() primitive.ObjectID
	SetID(primitive.ObjectID)
	// GetName() string
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

// 获取集合名词
func getCollectionName(m Model) (string, error) {
	/*
		var typ = reflect.TypeOf(o)

		for typ.Kind() == reflect.Ptr ||
			typ.Kind() == reflect.Slice ||
			typ.Kind() == reflect.Array {
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Interface {
			// log.Printf("reflect.Interface1: %v\n", typ)
			// typ = reflect.TypeOf(typ)
			// log.Printf("reflect.Interface2: %v\n", typ)
			// return GetCollectionName(typ)

			// var val = reflect.ValueOf(o)
		}
		var name = typ.Name()
		log.Printf("type name: %v\n", name)
		if name == "" {
			return ""
		}
	*/
	var typ = reflect.TypeOf(m)
	if typ.Kind() != reflect.Ptr {
		return "", errors.New("[mongodb]m should be a ptr")
	}
	typ = typ.Elem()
	var name = typ.Name()
	if name == "" {
		return "", errors.New("[mongodb]no fond a name")
	}

	return util.ToSnakeCase(name), nil
}

// GetByID 根据ID查询
func GetByID(db *mongo.Database, m Model) bool {
	if m.GetID().IsZero() {
		return false
	}

	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return false
	}

	if err := db.Collection(name).FindOne(nil, bson.D{{"_id", m.GetID()}}).Decode(m); err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

// UpdateByID 根据ID更新字段
func UpdateByID(db *mongo.Database, m Model, update interface{}) error {
	if m.GetID().IsZero() {
		return errors.New("需指定ID")
	}

	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).UpdateOne(nil, bson.D{{"_id", m.GetID()}}, update); err != nil {
		return err
	}
	return nil
}

// DeleteByID 根据ID删除
func DeleteByID(db *mongo.Database, m Model) error {
	if m.GetID().IsZero() {
		return errors.New("需指定ID")
	}

	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).DeleteOne(nil, bson.D{{"_id", m.GetID()}}); err != nil {
		return err
	}
	return nil
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

	name, err := getCollectionName(m)
	if err != nil {
		return err
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

	if err := db.Collection(name).FindOne(nil, filter, opts...).Decode(m); err == mongo.ErrNoDocuments {
		return err
	}

	return nil
}

// FindAll 查找一个结果
//
// filter 查询条件：bson.D{{"uid", 1}}
func FindAll(db *mongo.Database, m Model, results interface{}, filter interface{}, op *QueryOpts) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
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

	if cursor, err := db.Collection(name).Find(nil, filter, opts...); err != nil {
		return err
	} else {
		if err = cursor.All(nil, results); err != nil {
			return err
		}
	}

	return nil
}

func CreateUpdate(db *mongo.Database, m Model) error {
	if db == nil {
		db = GetDB()
	}
	name, err := getCollectionName(m)
	if err != nil {
		return err
	}
	if m.GetID().IsZero() { // create
		if _, err := db.Collection(name).InsertOne(nil, m); err != nil {
			return err
		}
	} else { // update
		if _, err := db.Collection(name).ReplaceOne(nil, bson.D{{"_id", m.GetID()}}, m); err != nil {
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
	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if res, err := db.Collection(name).InsertOne(nil, m); err != nil {
		return err
	} else {
		// log.Printf("结果：%+v\n", res)
		id, _ := res.InsertedID.(primitive.ObjectID)
		m.SetID(id)
	}
	return nil
}

// SaveAll 保存多个
func SaveAll(db *mongo.Database, m Model, data []interface{}) error {
	if len(data) == 0 {
		return errors.New("save empty data")
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if db == nil {
		db = GetDB()
	}

	if res, err := db.Collection(name).InsertMany(nil, data); err != nil {
		return err
	} else {
		for i, insertedID := range res.InsertedIDs {
			id, _ := insertedID.(primitive.ObjectID)
			if t, ok := data[i].(Model); ok {
				t.SetID(id)
			}
		}
	}
	return nil
}

// UpdateAll 指定条件更新所有
func UpdateAll(db *mongo.Database, m Model, filter interface{}, update interface{}) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).UpdateMany(nil, filter, update); err != nil {
		return err
	}
	return nil
}

// DeleteAll 指定条件删除所有
func DeleteAll(db *mongo.Database, m Model, filter interface{}) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).DeleteMany(nil, filter); err != nil {
		return err
	}
	return nil
}
