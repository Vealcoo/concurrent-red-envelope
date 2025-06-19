package mysql

import (
	"gorm.io/gorm"
)

func InsertTestData(db *gorm.DB) error {
	for i := int64(1); i <= 1000; i++ {
		tx := db.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		id, err := CreateCampaign(db)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = CreateRedEnvelopeByCampaignID(db, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}

	return nil
}
