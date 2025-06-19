package mysql

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func FindCampaignByID(db *gorm.DB, id int64) (*Campaign, error) {
	var campaign Campaign
	err := db.Model(&Campaign{}).Where("id = ?", id).First(&campaign).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func CreateCampaign(db *gorm.DB) (int64, error) {
	campaign := &Campaign{Status: 1}
	if err := db.Model(&Campaign{}).Create(campaign).Error; err != nil {
		return 0, err
	}
	return campaign.ID, nil
}

func FindCampaignByStatus(db *gorm.DB, status int32) ([]*Campaign, error) {
	campaigns := []*Campaign{}
	err := db.Model(&Campaign{}).Where("status = ?", status).Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func UpdateCampaignStatus(db *gorm.DB, id int64, statusFrom, statusTo int32) error {
	result := db.Model(&Campaign{}).Where("id = ?", id).Where("status = ?", statusFrom).Update("status", statusTo)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("update failed")
	}

	return nil
}

func CreateRedEnvelopeByCampaignID(db *gorm.DB, campaignID int64) error {
	envelopes := make([]RedEnvelope, 2000)
	for i := int64(0); i < 2000; i++ {
		envelopes[i] = RedEnvelope{
			CampaignID: campaignID,
			Status:     int32(0),
			UserID:     "",
			Balance:    rand.Int63n(100000) + 1,
			RewardSn:   fmt.Sprintf("SN-%d-%d", campaignID, i),
			RewardDesc: "Reward description",
			CreatedAt:  time.Now().Unix(),
			Creater:    "System",
			UpdatedAt:  time.Now().Unix(),
			Updater:    "System",
		}
	}

	result := db.CreateInBatches(&envelopes, 200)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateRedEnvelopeByID(db *gorm.DB, id, status int64, userID string) error {
	return db.Model(&RedEnvelope{}).Where("id = ?", id).Update("status", status).Update("user_id", userID).Error
}
