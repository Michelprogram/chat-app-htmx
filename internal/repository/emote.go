package repository

import (
	"context"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmoteRepository struct {
	collection string
	conn       *database.Mongo
}

func NewEmoteRepository(collection string, conn *database.Mongo) IEmoteRepository {
	return EmoteRepository{
		collection: collection,
		conn:       conn,
	}
}

func (em EmoteRepository) Insert(messageid string, emote *models.Emote) error {

	collection := em.conn.Database.Collection(em.collection)

	objectID, err := primitive.ObjectIDFromHex(messageid)

	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"emotes": emote}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (em EmoteRepository) AddCounter(emoteid string) (models.Emote, error) {

	var updatedDocument struct {
		Emotes []models.Emote `bson:"emotes"`
	}

	collection := em.conn.Database.Collection(em.collection)

	objectID, err := primitive.ObjectIDFromHex(emoteid)

	filter := bson.M{"emotes": bson.M{"$elemMatch": bson.M{"_id": objectID}}}

	update := bson.M{"$inc": bson.M{"emotes.$.counter": 1}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDocument)

	if err != nil {
		return models.Emote{}, err
	}

	for _, item := range updatedDocument.Emotes {
		if item.ID == objectID {
			return item, nil
		}
	}

	return models.Emote{}, nil
}
