package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IndexCreate(c *mongo.Collection) error {
	ctx := context.Background()
	cursor, err := c.Indexes().List(ctx, options.ListIndexes().SetMaxTime(30*time.Second))
	if err != nil {
		return err
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return err
	}

	campaignIDStatusUserID := false
	expireTime := false
	for _, v := range result {
		key := v["key"].(bson.M)
		if key["expire_time"] != nil {
			expireTime = true
		}
		if key["campaign_id"] != nil && key["status"] != nil && key["user_id"] != nil {
			campaignIDStatusUserID = true
		}
	}

	if !expireTime {
		_, err = c.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys:    bson.M{"expire_time": -1},
				Options: options.Index().SetExpireAfterSeconds(0),
			},
		)
		if err != nil {
			return err
		}
	}

	if !campaignIDStatusUserID {
		_, err = c.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys: bson.D{{"campaign_id", 1}, {"status", 1}, {"user_id", 1}},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
