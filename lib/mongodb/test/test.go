package main

import (
	"context"
	"encoding/json"
	"errors"
	"iris-project/lib/mongodb"
	"iris-project/lib/mongodb/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// saveOne()
	// saveAll()
	// updateAll()

	// getByID()
	// findOne()
	// findAll()
	// count()

	// createUpdate()

	// indexs()

	// transaction()

	// aggregate()
	// sum()

	// pageGroup()

}

func saveAll() {
	orders := make([]interface{}, 0)
	orders = append(orders, &model.Order{UID: 1, TradeNo: "2000", TotalPrice: mongodb.GetDecimal128(50.5), Discount: mongodb.GetDecimal128(12657), TotalNum: 1, OpLog: []string{"下单"}})
	orders = append(orders, &model.Order{UID: 1, TradeNo: "2001", TotalPrice: mongodb.GetDecimal128(100), Discount: mongodb.GetDecimal128(129.32), TotalNum: 1, OpLog: []string{"下单", "支付"}})
	orders = append(orders, &model.Order{UID: 2, TradeNo: "2002", TotalPrice: mongodb.GetDecimal128(200), Discount: mongodb.GetDecimal128("985522"), TotalNum: 2, OpLog: []string{"下单", "支付", "发货"}})
	orders = append(orders, &model.Order{UID: 2, TradeNo: "3000", TotalPrice: mongodb.GetDecimal128(50.5), Discount: mongodb.GetDecimal128("985522.45865"), TotalNum: 1, OpLog: []string{}})
	orders = append(orders, &model.Order{UID: 3, TradeNo: "3001", TotalPrice: mongodb.GetDecimal128(100), Discount: mongodb.GetDecimal128(9856587556.222567), TotalNum: 1})
	orders = append(orders, &model.Order{UID: 3, TradeNo: "3002", TotalPrice: mongodb.GetDecimal128(200), Discount: mongodb.GetDecimal128(-12536584.6585225), TotalNum: 2})
	err := mongodb.SaveAll(nil, nil, &model.Order{}, orders)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, order := range orders {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}

func page() {
	var (
		page int = 1
		size int = 2
	)
	searchList := &mongodb.SearchListData{
		M:      &model.Order{},
		Filter: bson.D{},
		Sort:   bson.D{{"created_at", -1}, {"_id", -1}},
		Page:   page,
		Size:   size,
	}

	var (
		orders []*model.Order
		total  int64
	)
	if err := searchList.GetList(&orders, &total); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(total)
	for _, order := range orders {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}

func pageGroup() {
	var (
		page int = 1
		size int = 2
	)
	searchList := &mongodb.SearchListData{
		M:      &model.Order{},
		Filter: bson.D{},
		Sort:   bson.D{{"created_at", -1}, {"_id", -1}},
		Page:   page,
		Size:   size,
		Group:  []string{"uid"},
		GroupField: bson.D{
			{"name", bson.D{{"$first", "$name"}}},
			{"all_price", bson.D{{"$sum", "$total_price"}}},
			{"all_num", bson.D{{"$sum", "$total_num"}}},
			{"count", bson.D{{"$sum", 1}}},
		},
	}

	type Res struct {
		UID      int     `bson:"_id" json:"uid"`
		Name     string  `bson:"name" json:"name"`
		AllPrice float64 `bson:"all_price" json:"all_price"`
		AllNum   int     `bson:"all_num" json:"all_num"`
		Count    int     `bson:"count" json:"count"`
	}

	var (
		// results []bson.M
		results []Res
		total   int64
	)
	if err := searchList.GetList(&results, &total); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(total)
	for _, result := range results {
		log.Printf("%+v\n", result)
	}
	log.Println("-------")
}

func aggregate() {
	var (
		results = make([]bson.M, 0)
		stages  = []bson.D{
			// {{"$match", bson.D{{"uid", 1}}}},

			{{"$group", bson.D{{"_id", nil}, {"all_price", bson.D{{"$sum", "$total_price"}}}}}},
		}
	)
	if err := mongodb.Aggregate(nil, nil, &model.Order{}, stages, &results); err != nil {
		log.Println(err.Error())
	}
	for _, result := range results {
		log.Printf("%v\n", result)
	}
	log.Println("-------")
}

func sum() {
	amount, err := mongodb.Sum(nil, nil, &model.Order{}, bson.D{{"uid", 2}}, "total_price")
	log.Printf("Amount: %v, Error: %v", amount, err)
	log.Println("-------")
}

func indexs() {
	if err := mongodb.Indexs(nil, nil, &model.Order{}, bson.D{{"trade_no", 1}}, true); err != nil {
		log.Println(err.Error())
	}
	log.Println("-------")
}

func getByID() {
	order := &model.Order{ID: mongodb.GetObjectIDFromStr("6478032e8d41c0b93d10cb84")}
	if mongodb.GetByID(nil, nil, order) {
		log.Printf("%+v\n", order)
	} else {
		log.Println("not found")
	}
	log.Println("-------")
}

func findOne() {
	var order = &model.Order{}
	if err := mongodb.FindOne(nil, nil, order, bson.D{{"uid", 1}}, &mongodb.QueryOpts{Sort: bson.D{{"total_price", -1}}}); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}

func findAll() {
	var orders []*model.Order
	if err := mongodb.FindAll(nil, nil, &model.Order{}, &orders, bson.D{{"uid", 1}}, nil); err != nil {
		log.Println(err.Error())
	} else {
		for _, order := range orders {
			log.Printf("%+v\n", order)
			bytes, _ := json.Marshal(order)
			log.Printf("%s\n", string(bytes))
		}
	}
	log.Println("-------")
}

func count() {
	count, err := mongodb.Count(nil, nil, &model.Order{}, bson.D{})
	log.Printf("Count: %d, Error: %v", count, err)
}

func createUpdate() {
	var order = &model.Order{
		UID: 4, TradeNo: "4001", TotalPrice: mongodb.GetDecimal128(196.5), TotalNum: 8,
	}
	if err := mongodb.CreateUpdate(nil, nil, order); err != nil {
		log.Printf("%v\n", err)
	} else {
		log.Printf("%+v\n", order)
	}

	// update
	order.TotalPrice = mongodb.GetDecimal128(268.25)
	order.TotalNum = 10
	if err := mongodb.CreateUpdate(nil, nil, order); err != nil {
		log.Printf("%v\n", err)
	} else {
		log.Printf("%+v\n", order)
	}

	log.Println("-------")
}

func saveOne() {
	order := &model.Order{UID: 1, TradeNo: "1005", Name: "何永亮", TotalPrice: mongodb.GetDecimal128(88.99), Discount: mongodb.GetDecimal128(0), TotalNum: 2, DeliverTime: primitive.NewDateTimeFromTime(time.Now())}
	err := mongodb.SaveOne(nil, nil, order)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%+v, %v\n", order, order.ID.Hex())
	log.Println("-------")
}

// func saveAll() {
// 	orders := make([]interface{}, 0)
// 	orders = append(orders, &model.Order{UID: 2, TradeNo: "2000", TotalPrice: mongodb.GetDecimal128(50.5), TotalNum: 1})
// 	orders = append(orders, &model.Order{UID: 2, TradeNo: "2001", TotalPrice: mongodb.GetDecimal128(100), TotalNum: 1})
// 	orders = append(orders, &model.Order{UID: 2, TradeNo: "2002", TotalPrice: mongodb.GetDecimal128(200), TotalNum: 2})
// 	orders = append(orders, &model.Order{UID: 3, TradeNo: "3000", TotalPrice: mongodb.GetDecimal128(50.5), TotalNum: 1})
// 	orders = append(orders, &model.Order{UID: 3, TradeNo: "3001", TotalPrice: mongodb.GetDecimal128(100), TotalNum: 1})
// 	orders = append(orders, &model.Order{UID: 3, TradeNo: "3002", TotalPrice: mongodb.GetDecimal128(200), TotalNum: 2})
// 	err := mongodb.SaveAll(nil, nil, &model.Order{}, orders)
// 	if err != nil {
// 		log.Printf("Error: %v\n", err)
// 	}
// 	for _, order := range orders {
// 		log.Printf("%+v\n", order)
// 	}
// 	log.Println("-------")
// }

func updateAll() {
	mongodb.UpdateAll(nil, nil, &model.Order{}, bson.D{{"uid", 2}}, bson.D{{"$set", bson.D{{"name", "黄晓明"}}}})
	mongodb.UpdateAll(nil, nil, &model.Order{}, bson.D{{"uid", 3}}, bson.D{{"$set", bson.D{{"name", "刘亦菲"}}}})
	mongodb.UpdateAll(nil, nil, &model.Order{}, bson.D{{"uid", 4}}, bson.D{{"$set", bson.D{{"name", "张三丰"}}}})
}

func transaction() {
	if err := mongodb.Transaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		order := &model.Order{ID: mongodb.GetObjectIDFromStr("6455f58d1776ee67a2b7439c")}
		if mongodb.GetByID(nil, ctx, order) {
			if err := mongodb.UpdateByID(nil, ctx, order, bson.D{{"$set", bson.D{{"total_price", 222}}}}); err != nil {
				return nil, err
			}
			if false {
				return nil, errors.New("error manual")
			}
			if err := mongodb.UpdateByID(nil, ctx, order, bson.D{{"$unset", bson.D{{"total_numa", ""}}}}); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}); err != nil {
		log.Println(err.Error())
	}
}
