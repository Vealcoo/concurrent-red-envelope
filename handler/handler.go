package handler

import (
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
	h := &Handler{
		MySQL:             mysql.MySQLClient,
		Redis:             redismodel.RedisClient,
		MongoDBCollection: mongodb.MongoClient.Database("test").Collection("red_envelope"),
	}

	return h
}
