package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Host             string `mapstructure:"HOST"`
	Port             int    `mapstructure:"PORT"`
	DBPort           int    `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	PeopleAPIBaseURL string `mapstructure:"PEOPLE_API_BASE_URL"`
}

func Load() (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the .env file: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
	return config, nil
}
