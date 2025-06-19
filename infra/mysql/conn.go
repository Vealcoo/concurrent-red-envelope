package mysql

import (
	"test/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Conn() (*gorm.DB, error) {
	dsn := config.Setting.GetString("mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)                // 最大開啟連線數
	sqlDB.SetMaxIdleConns(10)                 // 最大閒置連線數
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // 連線最長存活時間

	return db, nil
}
