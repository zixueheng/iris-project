/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-15 16:04:43
 * @LastEditTime: 2023-06-02 10:17:44
 */
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"iris-project/global"
	"iris-project/lib/util"
	"log"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
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

// GetDecimal128 获取 Decimal128，参数obj数字或数字字符串
func GetDecimal128(obj interface{}) primitive.Decimal128 {
	var (
		decimalStr string
		zero       = primitive.NewDecimal128(0, 0)
	)
	switch value := obj.(type) {
	case nil:
		return zero
	case int:
		decimalStr = decimal.NewFromInt(int64(value)).String()
	case int32:
		decimalStr = decimal.NewFromInt(int64(value)).String()
	case int64:
		decimalStr = decimal.NewFromInt(value).String()
	case uint:
		decimalStr = decimal.NewFromInt(int64(value)).String()
	case uint32:
		decimalStr = decimal.NewFromInt(int64(value)).String()
	case uint64:
		decimalStr = decimal.NewFromInt(int64(value)).String()
	case float32:
		decimalStr = decimal.NewFromFloat32(value).String()
	case float64:
		decimalStr = decimal.NewFromFloat(value).String()
	case string:
		decimalStr = value
	default:
		return zero
	}
	if decimal128, err := primitive.ParseDecimal128(decimalStr); err != nil {
		return zero
	} else {
		return decimal128
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
		return "", errors.New("[mongodb]not fond a name from m")
	}

	return util.ToSnakeCase(name), nil
}

// Indexs 创建索引，keys如：bson.D{{"title", 1}}
func Indexs(db *mongo.Database, ctx context.Context, m Model, keys bson.D, unique bool) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(unique),
	}
	if _, err := db.Collection(name).Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}
	return nil
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
	Skip       int64  // 跳过数量
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
//
// 特别注意：参数`m`和参数`results`切片元素必须是同一个类型
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
		if op.Skip != 0 {
			opts = append(opts, options.Find().SetSkip(op.Skip))
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

// Aggregate 聚合查询，可进行分组、合计、平均等聚合操作；results应该是个指针
//
// 聚合管道操作：https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/
//
// stages 示例：
//
//	[]bson.D{
//		bson.D{{"$match", bson.D{{"toppings", "milk foam"}}}} // 匹配字段条件
//		bson.D{{"$unset", bson.A{"_id", "category"}}} // 排除字段
//		bson.D{{"$sort", bson.D{{"price", 1}, {"toppings", 1}}}} // 排序
//		bson.D{{"$limit", 2}} // 限定数量
//	}
//
// 分组：
//
//	[]bson.D{
//	    bson.D{{"$group", bson.D{
//	        {"_id", "$category"}, // 按 category分组
//	        {"average_price", bson.D{{"$avg", "$price"}}}, // average_price字段统计 price的平均值
//	        {"type_total", bson.D{{"$sum", 1}}}, // type_total 合计数量
//	    }}}
//	}
//
// 查询数量：
//
//	[]bson.D{
//		bson.D{{"$match", bson.D{{"rating", bson.D{{"$gt", 5}}}}}},
//		bson.D{{"$count", "counted_documents"}},
//	}
//
// 结果：[{counted_documents 5}]
func Aggregate(db *mongo.Database, ctx context.Context, m Model, stages []bson.D, results interface{}) error {
	if db == nil {
		db = GetDB()
	}

	name, err := getCollectionName(m)
	if err != nil {
		return err
	}

	pipeline := mongo.Pipeline{}
	for _, stage := range stages {
		pipeline = append(pipeline, stage)
	}

	cursor, err := db.Collection(name).Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, results); err != nil {
		return err
	}

	return nil
}

// Sum 求和查询，filter查询条件，field求和字段
func Sum(db *mongo.Database, ctx context.Context, m Model, filter bson.D, field string) (amount interface{}, err error) {
	var (
		results = make([]bson.M, 0)
		stages  = []bson.D{
			{{"$match", filter}},
			{{"$group", bson.D{{"_id", nil}, {"amount", bson.D{{"$sum", "$" + field}}}}}},
		}
	)
	if err = Aggregate(nil, nil, m, stages, &results); err != nil {
		return
	}
	for _, result := range results {
		log.Printf("%v\n", result)
		if v, ok := result["amount"]; ok {
			return v, nil
		}
	}
	err = errors.New("not found field amount")
	return
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
		if len(f) == 0 { // 查询所有使用_id 字段作为索引加快查询速度
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

// CreateUpdate 新增或更新，ID存在更新不存在则新增
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
		if res, err := db.Collection(name).InsertOne(ctx, m); err != nil {
			return err
		} else {
			id, _ := res.InsertedID.(primitive.ObjectID)
			m.SetID(id)
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
// 参数`data`必须是Model类型的切片
//
// 特别注意：参数`m`和参数`data`切片元素必须是同一个类型
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

// SearchListData 通用列表查询条件
type SearchListData struct {
	Db               *mongo.Database
	Ctx              context.Context
	M                Model
	Filter           interface{}
	Projection, Sort bson.D
	Group            []string // 分组字段
	GroupField       bson.D   // 分组字段外额外字段，如已有字段或合计字段 bson.D{{"name", bson.D{{"$first", "$name"}}}, {"all_price", bson.D{{"$sum", "$total_price"}}}}
	Page, Size       int
}

// GetList 通用列表查询
func (s *SearchListData) GetList(results interface{}, total *int64) (err error) {
	var (
		skip = int64((s.Page - 1) * s.Size)
		size = int64(s.Size)
	)
	if len(s.Group) == 0 { // 普通查询
		opts := &QueryOpts{
			Projection: s.Projection,
			Sort:       s.Sort,
			Skip:       skip,
			Limit:      size,
		}
		// log.Printf("opts %+v\n", opts)
		if *total, err = Count(s.Db, s.Ctx, s.M, s.Filter); err != nil {
			return
		}
		if err = FindAll(s.Db, s.Ctx, s.M, results, s.Filter, opts); err != nil {
			return
		}
	} else { // 分组查询
		var (
			stages  = []bson.D{}
			stages2 = []bson.D{}
		)

		if s.Filter != nil {
			stages = append(stages, bson.D{{"$match", s.Filter}})
			stages2 = append(stages2, bson.D{{"$match", s.Filter}})
		}
		if s.Projection != nil {
			stages = append(stages, bson.D{{"$project", s.Projection}})
		}
		if len(s.Group) == 1 {
			var group = bson.D{{"_id", "$" + s.Group[0]}}
			for _, groupField := range s.GroupField {
				group = append(group, groupField)
			}
			stages = append(stages, bson.D{{"$group", group}})
			stages2 = append(stages2, bson.D{{"$group", group}})
		} else {
			var groupArr = []bson.E{}
			for _, s := range s.Group {
				groupArr = append(groupArr, bson.E{s, "$" + s})
			}
			var group = bson.D{{"_id", groupArr}}
			for _, groupField := range s.GroupField {
				group = append(group, groupField)
			}
			stages = append(stages, bson.D{{"$group", group}})
			stages2 = append(stages2, bson.D{{"$group", group}})
		}
		if s.Sort != nil {
			stages = append(stages, bson.D{{"$sort", s.Sort}})
		}
		stages = append(stages, bson.D{{"$skip", skip}})
		stages = append(stages, bson.D{{"$limit", size}})
		stages2 = append(stages2, bson.D{{"$count", "count"}})
		// log.Printf("stages %v\n", stages)
		// log.Printf("stages2 %v\n", stages2)
		if err = Aggregate(s.Db, s.Ctx, s.M, stages, results); err != nil {
			return
		}

		var countRes = make([]bson.M, 0)
		if err = Aggregate(s.Db, s.Ctx, s.M, stages2, &countRes); err != nil {
			return
		}

		for _, res := range countRes {
			if v, ok := res["count"]; ok {
				*total = int64(util.InterfaceToInt(v))
			}
		}
	}

	return
}
