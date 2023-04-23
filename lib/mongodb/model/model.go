package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderTest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID        int                `bson:"uid" json:"uid"`
	TradeNo    string             `bson:"trade_no" json:"trade_no"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	TotalNum   int                `bson:"total_num" json:"total_num"`
}

func (m *OrderTest) GetID() primitive.ObjectID {
	return m.ID
}
func (m *OrderTest) SetID(id primitive.ObjectID) {
	m.ID = id
}
