package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConversationBook stores the information about the conversation and the profile.
// each profile can behave differently on each conversation such as archiving the conversation or disposing
// the conversation. It will effects the only the conversation and the current profile not the other profiles
// on the conversation.
type ConversationBook struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	// Profile id (INDEX)
	Profile primitive.ObjectID `json:"profile" bson:"profile"`

	// Profile id (INDEX)
	Partner primitive.ObjectID `json:"partner" bson:"partner"`

	// Conversation id (INDEX)
	Conversation conversation `json:"conversation" bson:"conversation"`

	// IsArchive status of the conversation.
	//
	// If it is ARCHIVED conversation still continues to received the messages, but it will
	// not return on the conversation list if it is not specified.
	IsArchive bool `json:"isArchive" bson:"isArchive"`

	// Clear date of the conversation. It will not dispose the conversation conversation still
	// exists on the page but messages on the conversation will not be returned.
	ClearDate time.Time `json:"clearDate" bson:"clearDate"`

	// Dispose will remove the conversation from the conversation list.
	//
	// Disposed conversations still receive the messages but it will not be returned on the conversation list
	// even if it is specified on the conversation list.
	//
	// When the new messages comes to the disposed conversation then the conversation become active conversation
	// by changing the propery from true to false.
	Disposed bool `json:"disposed" bson:"disposed"`

	// Last message created at. When the new conversation created by default conversation id is used.
	LastMessageCreatedAt time.Time `json:"lastMessageCreatedAt" bson:"lastMessageCreatedAt"`

	// Status of the conversation book
	IsActive bool `json:"isActive" bson:"isActive"`

	LeaveDate time.Time `json:"leaveDate" bson:"leaveDate"`

	OldConversation *ConversationWithLookUp `json:"oldConversation" bson:"oldConversation"`

	// Blocked status for only SINGLE/PRIVATE conversations.
	BlockedStatus string `json:"blockedStatus" bson:"blockedStatus"`

	// Check if the conversation is muted or not. If it is muted then user will not be receive the push notification for that
	// conversation.
	Mute bool `json:"mute" bson:"mute"`
}

type conversation struct {
	// Conversation Id.
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	// Users who are in the conversation.
	Users []primitive.ObjectID `json:"users" bson:"users"`

	// Users who are in the conversation web socket ids.
	Sockets []string `json:"sockets" bson:"sockets"`

	// Type of the conversation. It can be GROUP or PRIVATE
	Type string `json:"type" bson:"type"`

	// Meta is the extra propery on the conversation. It depends application to application
	// it can be everything.
	Meta map[string]interface{} `json:"meta" bson:"meta"`

	// Last message id on the conversation. TODO: Add wht it is interface.
	LastMessage *primitive.ObjectID `json:"lastMessage" bson:"lastMessage"`

	// Last message created at. When the new conversation created by default conversation id is used.
	LastMessageCreatedAt time.Time `json:"lastMessageCreatedAt" bson:"lastMessageCreatedAt"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func (c *conversation) GetPartner(profileId *primitive.ObjectID) primitive.ObjectID {

	for _, u := range c.Users {

		if u.Hex() != profileId.Hex() {
			return u
		}
	}
	return *profileId
}

func (c *conversation) ConvertUsersIdsToString() []string {
	var ids []string
	for _, u := range c.Users {
		ids = append(ids, u.Hex())
	}
	return ids
}

// ConversationWithLookUp contains the conversation with the lookuped values.
type ConversationWithLookUp struct {
	// Conversation Id.
	ID primitive.ObjectID `json:"_id" bson:"_id"`

	// Users who are in the conversation.
	Users []Profile `json:"users" bson:"users"`

	// Users who are in the conversation web socket ids.
	Sockets []string `json:"sockets" bson:"sockets"`

	// Type of the conversation. It can be GROUP or PRIVATE
	Type string `json:"type" bson:"type"`

	// Meta is the extra propery on the conversation. It depends application to application
	// it can be everything.
	Meta *map[string]interface{} `json:"meta" bson:"meta"`

	// Last message id on the conversation. TODO: Add wht it is interface.
	LastMessage Message `json:"lastMessage" bson:"lastMessage"`

	// Last message created at. When the new conversation created by default conversation id is used.
	LastMessageCreatedAt time.Time `json:"lastMessageCreatedAt" bson:"lastMessageCreatedAt"`
}
