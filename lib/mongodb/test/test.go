/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-13 10:27:54
 * @LastEditTime: 2023-04-26 15:22:12
 */
package main

import (
	"context"
	"encoding/json"
	"iris-project/lib/mongodb"
	"iris-project/lib/mongodb/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	saveOne()
	saveAll()

	// getByID()
	// findOne()
	findAll()

	// transaction()
}

func getByID() {
	order := &model.Order{ID: mongodb.GetObjectIDFromStr("6440d1c4042ff0b52303effe")}
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
	if err := mongodb.FindAll(nil, nil, &model.Order{}, &orders, bson.D{{"uid", 3}}, nil); err != nil {
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

func saveOne() {
	order := &model.Order{UID: 1, TradeNo: "1000", TotalPrice: 88, TotalNum: 2}
	err := mongodb.SaveOne(nil, nil, order)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%+v, %v\n", order, order.ID.Hex())
	log.Println("-------")
}

func saveAll() {
	orders := make([]interface{}, 0)
	orders = append(orders, &model.Order{UID: 2, TradeNo: "2000", TotalPrice: 50.5, TotalNum: 1})
	orders = append(orders, &model.Order{UID: 2, TradeNo: "2001", TotalPrice: 100, TotalNum: 1})
	orders = append(orders, &model.Order{UID: 2, TradeNo: "2002", TotalPrice: 200, TotalNum: 2})
	orders = append(orders, &model.Order{UID: 3, TradeNo: "3000", TotalPrice: 50.5, TotalNum: 1})
	orders = append(orders, &model.Order{UID: 3, TradeNo: "3001", TotalPrice: 100, TotalNum: 1})
	orders = append(orders, &model.Order{UID: 3, TradeNo: "3002", TotalPrice: 200, TotalNum: 2})
	err := mongodb.SaveAll(nil, nil, &model.Order{}, orders)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, order := range orders {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}

func transaction() {
	if err := mongodb.Transaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		order := &model.Order{ID: mongodb.GetObjectIDFromStr("6440d1c4042ff0b52303effe")}
		if mongodb.GetByID(nil, ctx, order) {
			if err := mongodb.UpdateByID(nil, ctx, order, bson.D{{"$set", bson.D{{"total_price", 123}}}}); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}); err != nil {
		log.Println(err.Error())
	}
}
