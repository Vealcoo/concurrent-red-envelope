package cron

import (
	"log"
	"test/infra/mysql"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type server struct {
	cron  *cron.Cron
	mysql *gorm.DB
}

func New() *server {
	return &server{mysql: mysql.MySQLClient}
}

func (s *server) Run() {
	s.cron = cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	_, err := s.cron.AddFunc("*/1 * * * *", s.scanForDoneCampaign)
	if err != nil {
		log.Fatal(err)
	}

	s.cron.Start()
}

func (s *server) Close() {
	s.cron.Stop()
}
