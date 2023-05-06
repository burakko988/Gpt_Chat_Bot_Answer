package event

import (
	"chatbot/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Message struct {
		Id           string               `json:"_id"`
		Conversation string               `json:"conversation"`
		Sender       string               `json:"sender"`
		Content      string               `json:"content"`
		Type         string               `json:"type"`
		Caption      string               `json:"caption"`
		CreatedAt    time.Time            `json:"createdAt"`
		IsDeleted    bool                 `json:"isDeleted"`
		IsSeen       []primitive.ObjectID `json:"isSeen"`
		Reply        *model.ReplyMessage  `json:"reply"`
		SyncId       string               `json:"syncId"`
	}

	LastMessage struct {
		ID           primitive.ObjectID `json:"_id" bson:"_id"`
		Conversation primitive.ObjectID `json:"conversation"`
		Type         string             `json:"type"`
		Content      string             `json:"content"`
		Caption      string             `json:"caption"`
		Sender       string             `json:"sender"`
		CreatedAt    time.Time          `json:"createdAt"`
	}

	Conversation struct {
		ID                   primitive.ObjectID     `json:"_id" bson:"_id"`
		Type                 string                 `json:"type" bson:"type"`
		Meta                 map[string]interface{} `json:"meta" bson:"meta"`
		Users                []model.Profile        `json:"users"`
		LastMessage          LastMessage            `json:"lastMessage" `
		LastMessageCreatedAt time.Time              `json:"lastMessageCreatedAt" bson:"lastMessageCreatedAt"`
		IsBlocked            bool                   `json:"isBlocked" `
		IsBlockedBy          bool                   `json:"isBlockedBy" `
		IsActive             bool                   `json:"isActive"`
		IsArchive            bool                   `json:"isArchive" bson:"isArchive"`
	}

	MessageSocketEvent struct {
		Event          string             `json:"event"`
		ConversationId primitive.ObjectID `json:"conversationId"`
		Conversation   Conversation       `json:"conversation"`
		Message        Message            `json:"message"`
	}
)
