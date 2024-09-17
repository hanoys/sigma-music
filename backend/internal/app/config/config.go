package config

import (
	"github.com/JeremyLoy/config"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB struct {
		Type    string `yaml:"type"`
		Mongodb struct {
			Database string `yaml:"database"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			URL      string `yaml:"url"`
		} `yaml:"mongodb"`

		Postgres struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Name     string `yaml:"name"`
		} `yaml:"postgres"`
	} `yaml:"db"`

	JWT struct {
		AccessTokenExpTime  int64  `yaml:"access_expiration_time"`
		RefreshTokenExpTime int64  `yaml:"refresh_expiration_time"`
		SecretKey           string `yaml:"secret"`
	} `yaml:"jwt"`

	Redis struct {
		Host string `config:"REDIS_HOST"`
		Port string `config:"REDIS_PORT"`
	}

	Minio struct {
		Endpoint     string `config:"MINIO_ENDPOINT"`
		BucketName   string `config:"MINIO_BUCKET_NAME"`
		RootUser     string `config:"MINIO_ROOT_USER"`
		RootPassword string `config:"MINIO_ROOT_PASSWORD"`
	}

	Logger struct {
		LogLevel string `yaml:"level"`
	} `yaml:"log"`
}

func GetConfig(configPath string) (*Config, error) {
	var conf Config

	buf, err := os.ReadFile("./config/config.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &conf.DB)
	if err != nil {
		return nil, err
	}

	err = config.From(configPath).To(&conf.Redis)
	if err != nil {
		return nil, err
	}

	err = config.From(configPath).To(&conf.Minio)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
