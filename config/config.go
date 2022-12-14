package config

import (
	"github.com/BurntSushi/toml"
)

const configFilename = "config.toml"

type Config struct {
	Database struct {
		Host     string
		Username string
		Password string
		Database string
	}
	AWS struct {
		AccessKeyID     string
		SecretAccessKey string
		Region          string
		ArtifactsBucket string
	}
}

var Current Config

func Load() error {
	_, err := toml.DecodeFile(configFilename, &Current)
	if err != nil {
		return err
	}

	return nil
}
