package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server    Server    `mapstructure:"server"`
	MongoDB   MongoDB   `mapstructure:"mongodb"`
	RateLimit RateLimit `mapstructure:"rate_limit"`
}

type Server struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type MongoDB struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type Cleanup struct {
	Interval time.Duration `mapstructure:"interval"`
}

type RateLimit struct {
	Enabled bool  `mapstructure:"enabled"`
	Limit   int64 `mapstructure:"limit"`
	Burst   int   `mapstructure:"burst"`
}

// Load loads configuration from environment variables or config file
func Load() (*Config, error) {
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.read_timeout", 5*time.Second)
	viper.SetDefault("server.write_timeout", 10*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)

	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "linkbio")

	viper.SetDefault("cleanup.interval", 15*time.Minute)

	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.limit", 1000)
	viper.SetDefault("rate_limit.burst", 50)

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LINKBIO")

	// Config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
