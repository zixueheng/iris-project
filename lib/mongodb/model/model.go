package model

import (
	"iris-project/lib/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt  mongodb.LocalTime  `bson:"created_at" json:"created_at"`
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
