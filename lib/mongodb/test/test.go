/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-13 10:27:54
 * @LastEditTime: 2023-04-23 16:08:17
 */
package main

import (
	"iris-project/lib/mongodb"
	"iris-project/lib/mongodb/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	saveOne()
	saveAll()

	getByID()
	findOne()
	findAll()

}

func getByID() {
	id, _ := primitive.ObjectIDFromHex("6440d1c4042ff0b52303effe")
	order := &model.OrderTest{ID: id}
	if mongodb.GetByID(nil, order) {
		log.Printf("%+v\n", order)
	} else {
		log.Println("not found")
	}
	log.Println("-------")
}

func findOne() {
	var order = &model.OrderTest{}
	if err := mongodb.FindOne(nil, order, bson.D{{"uid", 1}}, &mongodb.QueryOpts{Sort: bson.D{{"total_price", -1}}}); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}

func findAll() {
	var orders []*model.OrderTest
	if err := mongodb.FindAll(nil, &model.OrderTest{}, &orders, bson.D{{"uid", 1}}, nil); err != nil {
		log.Println(err.Error())
	} else {
		for _, order := range orders {
			log.Printf("%+v\n", order)
		}
	}
	log.Println("-------")
}

func saveOne() {
	order := &model.OrderTest{UID: 3, TradeNo: "2000", TotalPrice: 300.33, TotalNum: 3}
	err := mongodb.SaveOne(nil, order)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%+v, %v\n", order, order.ID.Hex())
	log.Println("-------")
}

func saveAll() {
	orders := make([]interface{}, 0)
	orders = append(orders, &model.OrderTest{UID: 3, TradeNo: "2001", TotalPrice: 100, TotalNum: 1})
	orders = append(orders, &model.OrderTest{UID: 3, TradeNo: "2002", TotalPrice: 200, TotalNum: 2})
	err := mongodb.SaveAll(nil, &model.OrderTest{}, orders)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, order := range orders {
		log.Printf("%+v\n", order)
	}
	log.Println("-------")
}
