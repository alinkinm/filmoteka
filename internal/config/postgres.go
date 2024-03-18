package config

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	Name     string
	Password string
	Db       string
	Host     string
	Port     string
}

func GetDBConfig() (*DBConfig, error) {

	config := &DBConfig{}
	err := viper.UnmarshalKey("postgres", config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
