package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host      string `mapstructure:"HOST"`
	MongoURI  string `mapstructure:"MONGO_URI"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	GinMode   string `mapstructure:"GIN_MODE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")

	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
