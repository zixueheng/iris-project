/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:04:43
 * @LastEditTime: 2023-04-27 10:50:41
 */
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"iris-project/global"
	"iris-project/lib/util"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

// GetIDFromStr 获取ObjectID
func GetObjectIDFromStr(idStr string) primitive.ObjectID {
	if id, err := primitive.ObjectIDFromHex(idStr); err != nil {
		return primitive.NilObjectID
	} else {
		return id
	}
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
func GetByID(db *mongo.Database, ctx context.Context, m Model) bool {
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

	if err := db.Collection(name).FindOne(ctx, bson.D{{"_id", m.GetID()}}).Decode(m); err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

// UpdateByID 根据ID更新字段
//
// update 如 bson.D{{"$set", bson.D{{"field", "value"}}}}
func UpdateByID(db *mongo.Database, ctx context.Context, m Model, update interface{}) error {
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

	if _, err := db.Collection(name).UpdateOne(ctx, bson.D{{"_id", m.GetID()}}, update); err != nil {
		return err
	}
	return nil
}

// DeleteByID 根据ID删除
func DeleteByID(db *mongo.Database, ctx context.Context, m Model) error {
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

	if _, err := db.Collection(name).DeleteOne(ctx, bson.D{{"_id", m.GetID()}}); err != nil {
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
func FindOne(db *mongo.Database, ctx context.Context, m Model, filter interface{}, op *QueryOpts) error {
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

	if err := db.Collection(name).FindOne(ctx, filter, opts...).Decode(m); err == mongo.ErrNoDocuments {
		return err
	}

	return nil
}

// FindAll 查找一个结果
//
// filter 查询条件：bson.D{{"uid", 1}}
func FindAll(db *mongo.Database, ctx context.Context, m Model, results interface{}, filter interface{}, op *QueryOpts) error {
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

	if cursor, err := db.Collection(name).Find(ctx, filter, opts...); err != nil {
		return err
	} else {
		if err = cursor.All(ctx, results); err != nil {
			return err
		}
	}

	return nil
}

// Count 查询数量
func Count(db *mongo.Database, ctx context.Context, m Model, filter interface{}) (int64, error) {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return 0, err
	}
	var opts *options.CountOptions
	if f, ok := filter.(bson.D); ok {
		if len(f) == 0 { // 查询所有使用_id 字段作为索引加快查询熟读
			opts = options.Count().SetHint("_id_")
		}
	}

	return db.Collection(name).CountDocuments(ctx, filter, opts)
}

func setCreateAt(obj interface{}) error {
	// 获取对象的反射值
	objValue := reflect.ValueOf(obj)

	// 确保 obj 是指向结构体对象的指针类型
	if objValue.Kind() != reflect.Ptr || objValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("obj must be a pointer to struct")
	}

	// 获取结构体对象的反射值
	fieldValue := objValue.Elem().FieldByName("CreatedAt")

	// 确保字段存在并且可设置
	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil
	}

	now := NewLocalTime(time.Now())

	// 获取要设置的值的反射值
	valueOf := reflect.ValueOf(now)

	// 确保值类型与字段类型一致
	if fieldValue.Type() != valueOf.Type() {
		return fmt.Errorf("value type does not match field type")
	}

	// 设置字段值
	fieldValue.Set(valueOf)

	return nil
}

func CreateUpdate(db *mongo.Database, ctx context.Context, m Model) error {
	if db == nil {
		db = GetDB()
	}
	name, err := getCollectionName(m)
	if err != nil {
		return err
	}
	if m.GetID().IsZero() { // create
		setCreateAt(m)
		if _, err := db.Collection(name).InsertOne(ctx, m); err != nil {
			return err
		}
	} else { // update
		if _, err := db.Collection(name).ReplaceOne(ctx, bson.D{{"_id", m.GetID()}}, m); err != nil {
			return err
		}
	}
	return nil
}

// SaveOne 保存一个
func SaveOne(db *mongo.Database, ctx context.Context, m Model) error {
	if db == nil {
		db = GetDB()
	}
	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	setCreateAt(m)
	if res, err := db.Collection(name).InsertOne(ctx, m); err != nil {
		return err
	} else {
		// log.Printf("结果：%+v\n", res)
		id, _ := res.InsertedID.(primitive.ObjectID)
		m.SetID(id)
	}
	return nil
}

// SaveAll 保存多个
//
// data 必须是Model类型的切片
func SaveAll(db *mongo.Database, ctx context.Context, m Model, data []interface{}) error {
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
	for _, v := range data {
		setCreateAt(v)
	}

	if res, err := db.Collection(name).InsertMany(ctx, data); err != nil {
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
//
// update 如 bson.D{{"$set", bson.D{{"field", "value"}}}}
func UpdateAll(db *mongo.Database, ctx context.Context, m Model, filter interface{}, update interface{}) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).UpdateMany(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// DeleteAll 指定条件删除所有
func DeleteAll(db *mongo.Database, ctx context.Context, m Model, filter interface{}) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	if _, err := db.Collection(name).DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

// Transaction 执行事务操作
//
// 注意：事务只能在开启副本集的时候才能使用
func Transaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error)) error {
	// start-session
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := GetClient().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if _, err := session.WithTransaction(ctx, fn, txnOptions); err != nil {
		return err
	}
	// end-session

	return nil
}

/*
func T(){
	// start-session
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := GetClient().StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(nil)

	err = mongo.WithSession(nil, session, func(ctx mongo.SessionContext) error {
		if err = session.StartTransaction(txnOptions); err != nil {
			return err
		}

		docs := []interface{}{
			bson.D{{"title", "The Year of Magical Thinking"}, {"author", "Joan Didion"}},
			bson.D{{"title", "Play It As It Lays"}, {"author", "Joan Didion"}},
			bson.D{{"title", "The White Album"}, {"author", "Joan Didion"}},
		}
		result, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		// fmt.Println(result.InsertedIDs)
		return nil
	})
	if err != nil {
		if err := session.AbortTransaction(context.TODO()); err != nil {
			panic(err)
		}
		panic(err)
	}
}
*/
