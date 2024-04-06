package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Username       string             `bson:"username,omitempty" json:"username" form:"username"`
	Password       string             `bson:"password,omitempty" json:"password" form:"password"`
	ProfilePicture string             `bson:"profile_picture,omitempty" json:"profile_picture" form:"profile_picture"`
	LastConnection primitive.DateTime `bson:"last_connection,omitempty" json:"last_connection"`
	Token          string             `bson:"token,omitempty" json:"token"`
}

func NewUser(username, password string) *User {
	return &User{
		ID:             primitive.NewObjectID(),
		Username:       username,
		Password:       password,
		ProfilePicture: "https://rickandmortyapi.com/api/character/avatar/712.jpeg",
		LastConnection: primitive.DateTime(time.Now().Unix()),
		Token:          "",
	}
}
