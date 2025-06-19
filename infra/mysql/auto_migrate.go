package mysql

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&RedEnvelope{},
		&Campaign{},
	)
	if err != nil {
		return err
	}

	return nil
}
