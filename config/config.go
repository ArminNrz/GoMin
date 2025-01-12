package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	UseSSL    bool   `mapstructure:"useSSL"`
}

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access-key"`
		SecretKey string `yaml:"secret-key"`
		UseSSL    bool   `yaml:"use-ssl"`
	} `yaml:"minio"`
	API struct {
		Keys map[string]string `yaml:"api-key"`
	} `yaml:"api"`
	Jwt struct {
		Secret string `yaml:"secret"`
	}
}

var AppConfig Config

func LoadConfig() {

	profile := os.Getenv("PROFILE")
	if profile == "test" {
		profile = "-test"
	}

	configFile := fmt.Sprintf("config%s.yml", profile)

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	fmt.Println("Configuration loaded successfully")
}
