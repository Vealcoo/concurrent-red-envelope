package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GrabRedEnvelope(ctx context.Context, redisConn *redis.Client, campaignID int64) (*RedEnvelope, error) {
	queueKey := fmt.Sprintf("red_envelope_queue:%d", campaignID)
	data, err := redisConn.RPop(ctx, queueKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}

	envelope := &RedEnvelope{}
	if err := json.Unmarshal([]byte(data), &envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}
