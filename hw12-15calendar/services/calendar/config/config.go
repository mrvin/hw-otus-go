package config

import (
	"fmt"
	"io/ioutil"

	sqlstorage "github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage/sql"
	"gopkg.in/yaml.v2"
)

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggerConf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

type Config struct {
	InMem  bool              `yaml:"inmemory"`
	DB     sqlstorage.DBConf `yaml:"db"`
	HTTP   HTTPConf          `yaml:"http"`
	GRPC   GRPCConf          `yaml:"grpc"`
	Logger LoggerConf        `yaml:"logger"`
}

func Parse(configPath string) (*Config, error) {
	configYml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s error: %w", configPath, err)
	}

	var conf Config
	if err := yaml.Unmarshal(configYml, &conf); err != nil {
		return nil, fmt.Errorf("can't unmarshal %s: %w", configPath, err)
	}

	return &conf, nil
}
