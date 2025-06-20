package mongodb

import (
	"context"
	"log"
	"test/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var MongoClient *mongo.Client

func init() {
	clientOptions := options.Client().
		ApplyURI(config.Setting.GetString("mongodb.uri")).
		SetMaxPoolSize(1000).                // 最大連接數
		SetMinPoolSize(10).                  // 最小空閒連接數
		SetMaxConnIdleTime(30 * time.Second) // 空閒連接的最大存活時間
	if config.Setting.GetBool("mongodb.writeconcern.journal") {
		clientOptions.SetWriteConcern(writeconcern.Journaled()) // 寫關注
	}

	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	MongoClient = mongoClient

	err = IndexCreate(MongoClient.Database("test").Collection("red_envelope"))
	if err != nil {
		log.Fatal(err)
	}
}
