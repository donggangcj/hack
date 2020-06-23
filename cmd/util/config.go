package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	BlogConfig BlogConfig
}

type BlogConfig struct {
	BlogDir          string
	BlogSourceDir    string
	ImageRepoDir     string
	ImageRepoRootDir string
}

func BuildConfigFromFile(fPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(fPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
