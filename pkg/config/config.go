package config

import (
	"os"
	"path"
	"time"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/log"

	"gopkg.in/yaml.v3"
)

var configDir = os.Getenv("CONFIG_DIR")

const (
	configFileName = "config.yaml"
	secretFileName = "secret.yaml"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Server     `yaml:"server"`
	PostgreSQL `yaml:"postgresql"`
	S3         `yaml:"s3"`
}

type Server struct {
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type PostgreSQL struct {
	URI            string        `yaml:"uri"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

type S3 struct {
	Region   string `yaml:"region"`
	Endpoint string `yaml:"endpoint"`
}

func Load() (Config, error) {
	cfg := Config{}

	configFile, err := os.ReadFile(path.Join(configDir, configFileName))
	if err != nil {
		return cfg, err
	}
	secretFile, err := os.ReadFile(path.Join(configDir, secretFileName))
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(secretFile, &cfg)
	if err != nil {
		return cfg, err
	}

	log.DefaultLogger().Info().Msgf("loaded config: %+v", cfg)
	return cfg, nil
}
