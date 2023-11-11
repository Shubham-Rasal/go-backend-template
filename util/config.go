package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSource          string        `mapstructure:"DB_SOURCE"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {

	//check if ENV variable is set to ci or local
	//if ci then load config from environment variables
	//else load config from .env file
	viper.BindEnv("ENV")
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	env := viper.GetString("ENV")
	log.Println("ENV : ", env)
	log.Println("dbdriver : ", viper.GetString("DB_DRIVER"))
	if env == "ci" {
		viper.AutomaticEnv()
		err = viper.Unmarshal(&config)
		return
	}

	log.Println("loading config from .env file")
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
