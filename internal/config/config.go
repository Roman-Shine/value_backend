package config

import "time"

type (
	Config struct {
		Environment string
		HTTP        HTTPConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)
