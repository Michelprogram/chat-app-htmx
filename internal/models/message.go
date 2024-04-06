package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content,omitempty"`
	Author    primitive.ObjectID `bson:"author,omitempty"`
	Date      primitive.DateTime `bson:"date,omitempty"`
	IsDeleted bool               `bson:"isDeleted,omitempty"`

	Emotes []*Emote `bson:"emotes,omitempty"`
}

func NewMessage(content string, author primitive.ObjectID) *Message {
	return &Message{
		Content:   content,
		Author:    author,
		Date:      primitive.NewDateTimeFromTime(time.Now()),
		IsDeleted: false,
		Emotes:    make([]*Emote, 0),
	}
}

type EmbdedMessage struct {
	Message Message `bson:"message,omitempty"`
	User    User    `bson:"user,omitempty"`
}
