package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID                   primitive.ObjectID     `json:"_id" bson:"_id"`
	Sockets              []string               `json:"sockets" bson:"sockets"`
	Users                []primitive.ObjectID   `json:"users" bson:"users"`
	Type                 string                 `json:"type" bson:"type"`
	Meta                 map[string]interface{} `json:"meta" bson:"meta"`
	LastMessageCreatedAt time.Time              `json:"lastMessageCreatedAt" bson:"lastMessageCreatedAt"`
}

func (c *Conversation) GetBot(bids []primitive.ObjectID) *primitive.ObjectID {

	for _, u := range c.Users {

		for _, oi := range bids {

			if u.Hex() == oi.Hex() {
				return &oi
			}

		}
	}
	return nil
}
