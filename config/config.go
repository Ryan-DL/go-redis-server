package config

import (
	"os"
	"strconv"
)

type Config struct {
	RedisPassword *string
	RedisPort     *int
}

func LoadConfig() *Config {
	var cfg Config

	if redisPassword, exists := os.LookupEnv("REDIS_PASSWORD"); exists {
		cfg.RedisPassword = &redisPassword
	} else {
		cfg.RedisPassword = nil
	}

	if redisPortStr, exists := os.LookupEnv("REDIS_PORT"); exists {
		if redisPort, err := strconv.Atoi(redisPortStr); err == nil {
			cfg.RedisPort = &redisPort
		} else {
			cfg.RedisPort = nil
		}
	} else {
		cfg.RedisPort = nil
	}

	return &cfg
}
