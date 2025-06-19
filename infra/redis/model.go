package redis

type RedEnvelope struct {
	ID         int64
	CampaignID int64
	Status     int32
	UserID     string
	Amount     int64
}
