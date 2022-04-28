package config

import (
	"fmt"
	"io/ioutil"

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

type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type LoggerConf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

type Config struct {
	InMem  bool       `yaml:"inmemory"`
	DB     DBConf     `yaml:"db"`
	HTTP   HTTPConf   `yaml:"http"`
	GRPC   GRPCConf   `yaml:"grpc"`
	Logger LoggerConf `yaml:"logger"`
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
