package utils

import (
	"context"
	"log/slog"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Service struct {
		Port              string `toml:"port"`
		KeystoreFilePath  string `toml:"keystore_file_path"`
		KeystorePasswword string `toml:"keystore_password"`
	}
}

func NewConfig() *Config {
	return &Config{}
}

func InitConfig(configFile string, cfg *Config, logger *slog.Logger) (*Config, error) {

	_, err := toml.DecodeFile(configFile, &cfg)

	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "could not decode the config",
			slog.String("error", err.Error()))
		panic(err)
	}
	return cfg, nil
}
