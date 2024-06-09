package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reservation struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	RoomID   primitive.ObjectID `bson:"room_id" json:"room_id"`
	CheckIn  string             `bson:"check_in" json:"check_in"`
	CheckOut string             `bson:"check_out" json:"check_out"`
	Status   string             `bson:"status" json:"status"`
}
