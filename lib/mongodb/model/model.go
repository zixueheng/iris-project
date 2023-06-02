package model

import (
	"iris-project/lib/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Order struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	CreatedAt  mongodb.LocalTime  `bson:"created_at" json:"created_at"`
// 	UID        int                `bson:"uid" json:"uid"`
// 	Name       string             `bson:"name" json:"name"`
// 	TradeNo    string             `bson:"trade_no" json:"trade_no"`
// 	TotalPrice float64            `bson:"total_price" json:"total_price"`
// 	TotalNum   int                `bson:"total_num" json:"total_num"`
// }

// func (m *Order) GetID() primitive.ObjectID {
// 	return m.ID
// }
// func (m *Order) SetID(id primitive.ObjectID) {
// 	m.ID = id
// }

type Order struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CreatedAt   mongodb.LocalTime    `bson:"created_at" json:"created_at"`
	UID         int                  `bson:"uid" json:"uid"`
	Name        string               `bson:"name" json:"name"`
	TradeNo     string               `bson:"trade_no" json:"trade_no"`
	TotalPrice  primitive.Decimal128 `bson:"total_price" json:"total_price"` // Decimal128不可为nil
	Discount    primitive.Decimal128 `bson:"discount" json:"discount"`       // Decimal128不可为nil
	TotalNum    int                  `bson:"total_num" json:"total_num"`
	DeliverTime primitive.DateTime   `bson:"deliver_time" json:"deliver_time"`
	OpLog       []string             `bson:"op_log" json:"op_log"`
}

func (m *Order) GetID() primitive.ObjectID {
	return m.ID
}
func (m *Order) SetID(id primitive.ObjectID) {
	m.ID = id
}
