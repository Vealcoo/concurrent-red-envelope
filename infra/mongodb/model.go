package mongodb

import "time"

type RedEnvelope struct {
	ID         int64  `bson:"_id,omitempty"`
	CampaignID int64  `bson:"campaign_id"`
	Status     int32  `bson:"status"`
	UserID     string `bson:"user_id"`
	Amount     int64  `bson:"amount"`

	ExpireTime time.Time `bson:"expire_time"`
}
