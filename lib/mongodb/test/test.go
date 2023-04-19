/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-13 10:27:54
 * @LastEditTime: 2023-04-19 14:08:01
 */
package main

// "go.mongodb.org/mongo-driver/bson"
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

}
