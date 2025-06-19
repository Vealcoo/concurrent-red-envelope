package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var Setting *viper.Viper

func init() {
	Setting = setConfig("./conf.d", "setting", "yaml")
}

func setConfig(path, name, fileType string) *viper.Viper {
	conf := viper.New()
	conf.AddConfigPath(path)
	conf.SetConfigName(name)
	conf.SetConfigType(fileType)
	conf.AutomaticEnv()

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf.WatchConfig()

	go func() {
		timer := time.NewTimer(30 * time.Second)
		defer timer.Stop()
		for range timer.C {
			err := conf.ReadInConfig()
			if err != nil {
				log.Panicln(err)
			}
			timer.Reset(30 * time.Second)
		}
	}()

	return conf
}
