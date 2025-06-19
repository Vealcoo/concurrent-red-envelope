package cron

import (
	"test/infra/mysql"
	"time"
)

func (s *server) scanForDoneCampaign() {
	campaigns, err := mysql.FindCampaignByStatus(s.mysql, 2)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	for _, campaign := range campaigns {
		if campaign.ExpireTime < now {
			continue
		}

		err = mysql.UpdateCampaignStatus(s.mysql, campaign.ID, 2, 3)
		if err != nil {
			return
		}
	}
}
