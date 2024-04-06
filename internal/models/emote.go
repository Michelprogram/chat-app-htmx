package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Emote struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Link    string             `bson:"link,omitempty"`
	Counter int                `bson:"counter,omitempty"`
}

func NewEmote(link string) *Emote {
	return &Emote{
		Link:    link,
		ID:      primitive.NewObjectID(),
		Counter: 1,
	}
}
