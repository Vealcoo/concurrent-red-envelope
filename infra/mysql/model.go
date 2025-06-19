package mysql

type RedEnvelope struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`       // 紅包ID
	CampaignID int64  `gorm:"index:idx_campaign_status_user"` // 紅包活動ID
	Status     int32  `gorm:"index:idx_campaign_status_user"` // 紅包狀態 (0: 未領取, 1: 已領取)
	UserID     string `gorm:"index:idx_campaign_status_user"` // 領取用戶ID
	Balance    int64  // 紅包金額
	RewardSn   string // 獎項
	RewardDesc string // 獲獎訊息
	CreatedAt  int64  // 建立時間
	Creater    string // 建立者
	UpdatedAt  int64  // 修改時間
	Updater    string // 修改者
}

type Campaign struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	Status     int32 // 1: pending, 2:in progress, 3:done
	ExpireTime int64 // 到期時間, 開始後十分鐘過期
}
