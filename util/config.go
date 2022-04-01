package util

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/work/go/src/github.com/wham-invoice/wham-platform/secrets")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %v", err)
	}

}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}
