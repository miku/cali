package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Path string
	}
	Web struct {
		TemplatesDir string
		StaticDir    string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.path", "./cali.db")
	viper.SetDefault("web.templatesdir", "./web/templates")
	viper.SetDefault("web.staticdir", "./web/static")

	// Look for config in standard locations
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cali/")
	viper.AddConfigPath("$HOME/.config/cali/")
	viper.AddConfigPath(".")

	var config Config

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Ensure database path is absolute
	if !filepath.IsAbs(config.Database.Path) {
		absPath, err := filepath.Abs(config.Database.Path)
		if err != nil {
			return nil, err
		}
		config.Database.Path = absPath
	}

	return &config, nil
}
