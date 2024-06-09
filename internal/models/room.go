package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number string             `bson:"number" json:"number"`
	Type   string             `bson:"type" json:"type"`
	Price  float64            `bson:"price" json:"price"`
	Status string             `bson:"status" json:"status"`
}
