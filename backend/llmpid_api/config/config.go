package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all application settings
type Config struct {
	Host       HostConfiguration
	Database   DatabaseConfiguration
	Classifier ClassifierConfiguration
}

type HostConfiguration struct {
	Port        string
	Environment string
	LogsDirPath string
}

type DatabaseConfiguration struct {
	Host     string
	Port     string
	Name     string
	User     string // Loaded from ENV
	Password string // Loaded from ENV
}

type ClassifierConfiguration struct {
	ClassifierAPIPath string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Setting the default logs directory to <parent_dir>/logs in case it was not defined in the configuration.
	viper.SetDefault("LogsDirPath", "./logs")

	// Allow environment variables to be loaded.
	viper.AutomaticEnv()

	// Read the core config file.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, relying on environment variables.")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Override database credentials from environment variables.
	cfg.Database.User = viper.GetString("DB_USER")
	cfg.Database.Password = viper.GetString("DB_PASSWORD")

	return &cfg
}
