package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SendMessageToConversation struct {
		ConversationId string `json:"conversationId" bson:"conversationId"`
		Content        string `json:"content" bson:"content"`
		ContentType    string `json:"contentType" bson:"contentType"`
		Caption        string `json:"caption" bson:"caption"`
		Reply          string `json:"reply" bson:"reply"`
		SyncId         string `json:"syncId"`
		IsBot          string `json:"isBot"`
	}

	ReplyMessage struct {
		Id      primitive.ObjectID `json:"_id"`
		Content string             `json:"content"`
		Caption string             `json:"caption"`
		Type    string             `json:"type"`
		Sender  primitive.ObjectID `json:"sender" bson:"sender"`
	}

	WebSocketBody struct {
		Action string                    `json:"action" bson:"action"`
		Data   SendMessageToConversation `json:"data" bson:"data"`
	}

	Message struct {
		ID           primitive.ObjectID   `json:"_id" bson:"_id"`
		Sender       primitive.ObjectID   `json:"sender" bson:"sender"`
		Conversation primitive.ObjectID   `json:"conversation" bson:"conversation"`
		IsSeen       []primitive.ObjectID `json:"isSeen" bson:"isSeen"`
		Type         string               `json:"type" bson:"type"`
		Content      string               `json:"content" bson:"content"`
		Caption      string               `json:"caption" bson:"caption"`
		CreatedAt    time.Time            `json:"createdAt" bson:"createdAt"`
		IsDeleted    []primitive.ObjectID `json:"isDeleted" bson:"isDeleted"`
		Reply        interface{}          `json:"reply" bson:"reply"`
		IsCleared    []primitive.ObjectID `json:"isCleared" bson:"isCleared"`
	}
)
