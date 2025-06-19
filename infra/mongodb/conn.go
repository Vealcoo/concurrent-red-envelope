package mongodb

import (
	"context"
	"test/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Conn() (*mongo.Database, error) {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Setting.GetString("mongodb.uri")))
	if err != nil {
		return nil, err
	}
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database("test"), nil
}
