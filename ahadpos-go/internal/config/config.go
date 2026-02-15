package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	JWT       JWTConfig       `yaml:"jwt"`
	Log       LogConfig       `yaml:"log"`
	CORS      CORSConfig      `yaml:"cors"`
	Pagination PaginationConfig `yaml:"pagination"`
	Auth      AuthConfig      `yaml:"auth"`
}

type ServerConfig struct {
	Port            int           `yaml:"port"`
	Mode            string        `yaml:"mode"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Path            string `yaml:"path"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type CORSConfig struct {
	Enabled        bool     `yaml:"enabled"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	MaxAge         int      `yaml:"max_age"`
}

type PaginationConfig struct {
	DefaultPage     int `yaml:"default_page"`
	DefaultPageSize int `yaml:"default_page_size"`
	MaxPageSize     int `yaml:"max_page_size"`
}

type AuthConfig struct {
	BcryptCost int `yaml:"bcrypt_cost"`
}

var AppConfig *Config

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default values
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10
	}
	if config.JWT.Expiration == 0 {
		config.JWT.Expiration = 24
	}
	if config.Auth.BcryptCost == 0 {
		config.Auth.BcryptCost = 10
	}
	if config.Pagination.DefaultPage == 0 {
		config.Pagination.DefaultPage = 1
	}
	if config.Pagination.DefaultPageSize == 0 {
		config.Pagination.DefaultPageSize = 20
	}
	if config.Pagination.MaxPageSize == 0 {
		config.Pagination.MaxPageSize = 100
	}

	AppConfig = &config
	return &config, nil
}
