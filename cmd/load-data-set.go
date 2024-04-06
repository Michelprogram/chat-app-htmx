package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"math/rand"
	"time"
)

func LoadData() {

	users := LoadUser(5)

	LoadMessage(users, 50)

}

func LoadUser(n int) []*models.User {

	users := make([]*models.User, n)

	for i := 0; i < n; i++ {

		user, _ := models.InsertUser(faker.Username(), "dorian")

		users[i] = user

	}

	return users
}

func LoadMessage(users []*models.User, n int) []*models.Message {

	messages := make([]*models.Message, n)

	size := len(users)

	for i := 0; i < n; i++ {

		message, _ := InsertMessage(faker.Paragraph(), users[rand.Intn(size)].ID, randate())

		messages[i] = message

	}

	return messages

}

func InsertMessage(content string, author primitive.ObjectID, date time.Time) (*models.Message, error) {

	message := &models.Message{
		Content:   content,
		Author:    author,
		Date:      primitive.NewDateTimeFromTime(date),
		IsDeleted: false,
	}

	collection := database.Database.Collection(models.MessageCollection)

	result, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return nil, err
	}

	message.ID = result.InsertedID.(primitive.ObjectID)

	return message, nil
}

func randate() time.Time {
	min := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
