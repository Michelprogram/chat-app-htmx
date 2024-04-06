package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/event"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
	Debug    bool
}

func NewMongo(debug bool, name string) (*Mongo, error) {

	client, database, err := conn(debug, name)

	if err != nil {
		return nil, err
	}

	return &Mongo{
		Client:   client,
		Database: database,
		Debug:    debug,
	}, nil

}

func conn(debug bool, name string) (*mongo.Client, *mongo.Database, error) {

	url := os.Getenv("URL_MONGODB")

	if url == "" {
		return nil, nil, errors.New("can't find URL_MONGODB in .env")
	}

	connector := options.Client().ApplyURI(url)

	if debug {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		}

		connector.SetMonitor(cmdMonitor)
	}

	client, err := mongo.Connect(context.TODO(), connector)

	ctxTimeout, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	err = client.Ping(ctxTimeout, nil)

	if err != nil {
		return nil, nil, err
	}

	return client, client.Database(name), nil
}

func (m Mongo) Close() error {
	return m.Client.Disconnect(context.Background())
}
