package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	collection string
	conn       *database.Mongo
}

func NewUserRepository(collection string, conn *database.Mongo) IUserRepository {
	return UserRepository{
		collection: collection,
		conn:       conn,
	}
}

func hexIdsToPrimitiveObject(ids ...string) []primitive.ObjectID {

	idsObj := make([]primitive.ObjectID, 0, len(ids))

	for _, id := range ids {
		objId, _ := primitive.ObjectIDFromHex(id)
		idsObj = append(idsObj, objId)
	}

	return idsObj
}

func (r UserRepository) Insert(username, password string) (*models.User, error) {

	collection := r.conn.Database.Collection(r.collection)

	hash := sha256.Sum256([]byte(password))

	password = fmt.Sprintf("%x", hash[:])

	user := models.NewUser(username, password)

	id, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		return nil, err
	}

	user.ID = id.InsertedID.(primitive.ObjectID)

	return user, nil

}

func (r UserRepository) Find(username, password string) *models.User {

	var result models.User

	collection := r.conn.Database.Collection(r.collection)

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	filter := bson.M{"username": username, "password": hash}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil
	}

	return &result

}

func (r UserRepository) FindById(id string) *models.User {

	var result models.User

	collection := r.conn.Database.Collection(r.collection)

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil
	}

	filter := bson.M{"_id": objId}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil
	}

	return &result

}

func (r UserRepository) FindByUsername(username string) *models.User {

	var result models.User

	collection := r.conn.Database.Collection(r.collection)

	filter := bson.M{"username": username}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil
	}

	return &result

}

func (r UserRepository) FindAndUpdate(id string, user models.User) error {

	var update bson.M

	collection := r.conn.Database.Collection(r.collection)

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": objId}

	if user.Password != "" {
		update = bson.M{"$set": bson.M{
			"username":        user.Username,
			"profile_picture": user.ProfilePicture,
		},
		}
	} else {
		update = bson.M{"$set": user}
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil

}

func (r UserRepository) UpdateToken(id primitive.ObjectID, token string) error {
	var update bson.M

	collection := r.conn.Database.Collection(r.collection)

	filter := bson.M{"_id": id}

	update = bson.M{"$set": bson.M{"token": token}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) UsersnamesByIDs(ids []string) ([]models.User, error) {

	var users []models.User

	idsObj := hexIdsToPrimitiveObject(ids...)

	collection := r.conn.Database.Collection(r.collection)

	filter := bson.M{"_id": bson.M{"$in": idsObj}}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &users)

	if err != nil {
		return nil, err
	}

	return users, nil

}
