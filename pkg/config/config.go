package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

const configPath = "CONFIG_PATH"

// Config is where the global configuration is stored.
type Config struct {
	Port string `mapstructure:"PORT"`
	DB   dB     `mapstructure:",squash"`
	JWT  jwt    `mapstructure:",squash"`
}

type dB struct {
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
}

type jwt struct {
	Expiration int32  `mapstructure:"JWT_EXPIRATION_SECS"`
	Key        string `mapstructure:"JWT_KEY"`
}

// LoadConfig gets the configuration in from .env files and stores the in Config struct.
func LoadConfig() *Config {
	var c Config
	var once sync.Once

	once.Do(func() {
		cPath := os.Getenv(configPath)
		if cPath == "" {
			log.Fatalln("CONFIG_PATH is not set")
		}

		viper.AddConfigPath(cPath + "/pkg/config/env")
		viper.SetConfigName("dev")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalln("Could not read configuration:", err.Error())
		}

		err = viper.Unmarshal(&c)
		if err != nil {
			log.Fatalln("Could not unmarshal to Config struct:", err.Error())
		}
	})

	return &c
}
