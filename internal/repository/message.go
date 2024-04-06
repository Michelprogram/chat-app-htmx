package repository

import (
	"context"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var BUCKET_MESSAGE_SIZE int64 = 10

type MessageRepository struct {
	collection string
	conn       *database.Mongo
}

func NewMessageRepository(collection string, conn *database.Mongo) IMessageRepository {
	return MessageRepository{
		collection: collection,
		conn:       conn,
	}
}

func (mr MessageRepository) Insert(content string, author primitive.ObjectID) (*models.Message, error) {

	message := models.NewMessage(content, author)

	collection := mr.conn.Database.Collection(mr.collection)

	result, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return nil, err
	}

	message.ID = result.InsertedID.(primitive.ObjectID)

	return message, nil
}

func (mr MessageRepository) Delete(id, userID string) (*models.Message, error) {

	var result models.Message

	collection := mr.conn.Database.Collection(mr.collection)

	messageObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": messageObjectID, "author": userObjectID}

	update := bson.M{"$set": bson.M{"isDeleted": true}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)

	log.Println(result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func (mr MessageRepository) GetLast(n int64, from primitive.ObjectID) ([]models.EmbdedMessage, error) {

	var result []models.EmbdedMessage

	var message models.Message

	collection := mr.conn.Database.Collection(mr.collection)

	filter := bson.M{"_id": from}

	err := collection.FindOne(context.TODO(), filter).Decode(&message)

	if err != nil {
		return nil, err
	}

	matchStage := bson.D{{"$match", bson.D{{"date", bson.D{{"$lte", message.Date}}}}}}

	lookupStage := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "author"},
		{"foreignField", "_id"},
		{"as", "user"},
	}}}

	projectStage := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"message", bson.D{
			{"content", "$content"},
			{"author", "$author"},
			{"date", "$date"},
			{"_id", "$_id"},
			{"emotes", "$emotes"},
			{"isDeleted", "$isDeleted"},
		}},
		{"user", bson.D{{"$arrayElemAt", bson.A{"$user", 0}}}},
	}}}

	sortStage := bson.D{{"$sort", bson.D{{"message.date", -1}}}}
	limitStage := bson.D{{"$limit", n}}

	aggregate := bson.A{matchStage, lookupStage, projectStage, sortStage, limitStage}

	cursor, err := collection.Aggregate(context.TODO(), aggregate)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (mr MessageRepository) GetLastInsert() (models.Message, error) {

	var message models.Message

	collection := mr.conn.Database.Collection(mr.collection)

	opts := options.FindOne().SetSort(bson.M{"_id": -1})

	err := collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&message)

	if err != nil {
		return models.Message{}, err
	}

	return message, nil

}

func (mr MessageRepository) GetLasts(n int64) ([]models.EmbdedMessage, error) {

	var result []models.EmbdedMessage

	collection := mr.conn.Database.Collection(mr.collection)

	aggregate := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "author"},
					{"foreignField", "_id"},
					{"as", "user"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"message",
						bson.D{
							{"content", "$content"},
							{"author", "$author"},
							{"date", "$date"},
							{"_id", "$_id"},
							{"emotes", "$emotes"},
							{"isDeleted", "$isDeleted"},
						},
					},
					{"user",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$user",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"message.date", -1}}}},
		bson.D{{"$limit", n}},
	}

	cursor, err := collection.Aggregate(context.TODO(), aggregate)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil

}
