package handler

import (
	"fmt"
	"test/infra/mongodb"
	"test/infra/mysql"
	redismodel "test/infra/redis"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Handler struct {
	MySQL             *gorm.DB
	MongoDBCollection *mongo.Collection
	Redis             *redis.Client
}

func NewHandler() *Handler {
	mySQL, err := mysql.Conn()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	redisClient, err := redismodel.Conn()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	mongoDatabase, err := mongodb.Conn()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to mongodb: %v", err))
	}

	h := &Handler{
		MySQL:             mySQL,
		Redis:             redisClient,
		MongoDBCollection: mongoDatabase.Collection("red_envelope"),
	}

	return h
}
