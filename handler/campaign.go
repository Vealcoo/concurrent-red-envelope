package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"test/infra/mongodb"
	"test/infra/mysql"
	"test/infra/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) InsertTestData(c *gin.Context) {
	err := mysql.AutoMigrate(h.MySQL)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = mysql.InsertTestData(h.MySQL)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nil)
}

type CreateCampaignResponse struct {
	CampaignID int64 `json:"campaign_id"`
}

func (h *Handler) CreateCampaign(c *gin.Context) {
	err := h.MySQL.AutoMigrate(&mysql.RedEnvelope{}, &mysql.Campaign{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := mysql.CreateCampaign(h.MySQL)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = mysql.CreateRedEnvelopeByCampaignID(h.MySQL, id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp := CreateCampaignResponse{
		CampaignID: id,
	}
	c.JSON(200, resp)
}

type StartCampaignRequest struct {
	CacheMode string `json:"cache_mode" binding:"required"` // mongodb or redis
}

type StartCampaignResponse struct {
	CampaignID string `json:"campaign_id"`
}

func (h *Handler) StartCampaign(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	if campaignID == "" {
		c.JSON(400, gin.H{"error": "campaign_id is required"})
		return
	}

	var err error
	var campaignIDInt64 int64

	campaignIDInt64, err = strconv.ParseInt(campaignID, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid campaign_id"})
		return
	}

	var req StartCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	envelopes := []*mysql.RedEnvelope{}
	err = h.MySQL.Model(&mysql.RedEnvelope{}).Where("campaign_id = ?", campaignIDInt64).Find(&envelopes).Error
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(envelopes) == 0 {
		c.JSON(400, gin.H{"error": "campaign not found"})
		return
	}

	tx := h.MySQL.Begin()
	if tx.Error != nil {
		c.JSON(400, gin.H{"error": tx.Error.Error()})
		return
	}
	err = mysql.UpdateCampaignStatus(tx, campaignIDInt64, 1, 2)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	switch req.CacheMode {
	case "mongodb":
		envelopesForMongo := make([]any, len(envelopes))
		for i, envelope := range envelopes {
			envelopesForMongo[i] = &mongodb.RedEnvelope{
				ID:         envelope.ID,
				CampaignID: envelope.CampaignID,
				Status:     envelope.Status,
				UserID:     envelope.UserID,
				Amount:     int64(envelope.Balance),
				ExpireTime: time.Now().Add(time.Minute * 10),
			}
		}
		_, err = h.MongoDBCollection.InsertMany(ctx, envelopesForMongo)
		if err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

	case "redis":
		queueKey := fmt.Sprintf("red_envelope_queue:%s", campaignID)
		pipe := h.Redis.Pipeline()

		for _, envelope := range envelopes {
			data, err := json.Marshal(&redis.RedEnvelope{
				ID:         envelope.ID,
				CampaignID: envelope.CampaignID,
				Status:     envelope.Status,
				UserID:     envelope.UserID,
				Amount:     int64(envelope.Balance),
			})
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			pipe.LPush(ctx, queueKey, data)
		}
		pipe.Expire(ctx, queueKey, 10*time.Minute)
		_, err = pipe.Exec(ctx)
		if err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

	default:
		tx.Rollback()
		c.JSON(400, gin.H{"error": "invalid cache_mode"})
		return
	}
	tx.Commit()

	resp := CreateCampaignResponse{
		CampaignID: campaignIDInt64,
	}
	c.JSON(200, resp)
}

type ClaimRedEnvelopeRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	CacheMode string `json:"cache_mode" binding:"required"` // mongodb or redis
}

type ClaimRedEnvelopeResponse struct {
	RedEnvelopeID int64           `json:"red_envelope_id"`
	Amount        decimal.Decimal `json:"amount"`
}

func (h *Handler) ClaimRedEnvelope(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	if campaignID == "" {
		c.JSON(400, gin.H{"error": "campaign_id is required"})
		return
	}
	campaignIDInt64, err := strconv.ParseInt(campaignID, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid campaign_id"})
		return
	}

	var req ClaimRedEnvelopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	redEnvelopeID := int64(0)
	redEnvelopeAmount := decimal.Zero
	ctx := c.Request.Context()
	switch req.CacheMode {
	case "mongodb":
		res, err := mongodb.FindAndUpdateStatus(ctx, h.MongoDBCollection, campaignIDInt64, req.UserID)
		if err != nil && err != mongo.ErrNoDocuments {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if res == nil {
			c.JSON(400, gin.H{"error": "campaign finished"})
			return
		}

		err = mysql.UpdateRedEnvelopeByID(h.MySQL, res.ID, int64(res.Status), req.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		redEnvelopeID = res.ID
		redEnvelopeAmount = decimal.NewFromInt(res.Amount)

	case "redis":
		res, err := redis.GrabRedEnvelope(ctx, h.Redis, campaignIDInt64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if res == nil {
			c.JSON(400, gin.H{"error": "campaign finished"})
			return
		}

		err = mysql.UpdateRedEnvelopeByID(h.MySQL, res.ID, int64(res.Status)+1, req.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		redEnvelopeID = res.ID
		redEnvelopeAmount = decimal.NewFromInt(res.Amount)

	default:
		c.JSON(400, gin.H{"error": "invalid cache_mode"})
		return
	}

	resp := ClaimRedEnvelopeResponse{
		RedEnvelopeID: redEnvelopeID,
		Amount:        redEnvelopeAmount,
	}

	c.JSON(200, resp)
}
