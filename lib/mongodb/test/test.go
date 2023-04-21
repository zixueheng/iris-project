/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-13 10:27:54
 * @LastEditTime: 2023-04-20 15:06:54
 */
package main

import (
	"iris-project/lib/mongodb"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// "go.mongodb.org/mongo-driver/mongo"
// "go.mongodb.org/mongo-driver/mongo/options"
// "go.mongodb.org/mongo-driver/mongo/readpref"

func main() {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:123@localhost:27017/"))
		if err != nil {
			panic(err)
		}

		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
			cancel()
		}()

		// ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		// defer cancel()
		err = client.Ping(ctx, readpref.Primary())

		collection := client.Database("test").Collection("numbers")

		// ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()
		res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
		id := res.InsertedID
		log.Println(id)
	*/

	// saveOne()
	// saveAll()

	// id, _ := primitive.ObjectIDFromHex("643fa8337d71693bcd51cd8a")
	// order := &Order{ID: id}
	// if mongodb.GetByID(nil, order) {
	// 	log.Printf("%+v\n", order)
	// } else {
	// 	log.Println("not found")
	// }

	var order = &Order{}
	if err := mongodb.FindOne(nil, order, bson.D{{"uid", 1}}, &mongodb.QueryOpts{Sort: bson.D{{"total_price", -1}}}); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("%+v\n", order)
	}

}

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID        int                `bson:"uid" json:"uid"`
	TradeNo    string             `bson:"trade_no" json:"trade_no"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	TotalNum   int                `bson:"total_num" json:"total_num"`
}

func (m *Order) GetID() primitive.ObjectID {
	return m.ID
}
func (m *Order) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *Order) GetName() string {
	return "order"
}

func saveOne() {
	order := &Order{UID: 1, TradeNo: "1000", TotalPrice: 300.33, TotalNum: 3}
	err := mongodb.SaveOne(nil, order)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("%+v, %v\n\n", order, order.ID.Hex())
}

func saveAll() {
	orders := make([]interface{}, 0)
	orders = append(orders, &Order{UID: 1, TradeNo: "1001", TotalPrice: 100, TotalNum: 1})
	orders = append(orders, &Order{UID: 2, TradeNo: "1002", TotalPrice: 200, TotalNum: 2})
	err := mongodb.SaveAll(nil, orders)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, order := range orders {
		log.Printf("%+v\n\n", order.(*Order))
	}

}
