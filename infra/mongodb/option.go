package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAndUpdateRedEnvelopeStatus(ctx context.Context, c *mongo.Collection, campaignID int64, userID string) (*RedEnvelope, error) {
	res := &RedEnvelope{}
	err := c.FindOneAndUpdate(ctx, bson.M{
		"campaign_id": campaignID,
		"status":      0,
		"user_id": bson.M{
			"$in": []string{userID, ""},
		},
	}, bson.M{
		"$set": bson.M{
			"status":  1,
			"user_id": userID,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}
