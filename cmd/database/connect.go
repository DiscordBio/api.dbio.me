package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Disconnect(client *mongo.Client) error {
	err := client.Disconnect(nil)
	if err != nil {
		return err
	}
	return nil
}
