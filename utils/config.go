package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort       string `mapstructure:"APP_PORT"`
	RuntimeSetup  string `mapstructure:"RUNTIME_SETUP"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)

	return
}
