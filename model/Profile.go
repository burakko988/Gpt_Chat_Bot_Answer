package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	ProfileImage string             `json:"profileImage" bson:"profileImage"`
	Meta         meta               `json:"meta" bson:"meta"`
	Sockets      []string           `json:"sockets" bson:"sockets"`
}

type meta struct {
	Gender string `json:"sex" bson:"sex"`
	Age    int    `json:"age" bson:"age"`
	Name   string `json:"name" bson:"name"`
}
