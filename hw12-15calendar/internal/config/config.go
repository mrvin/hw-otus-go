package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Parse(configPath string, conf interface{}) error {
	configYml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading %s error: %w", configPath, err)
	}

	if err := yaml.Unmarshal(configYml, conf); err != nil {
		return fmt.Errorf("can't unmarshal %s: %w", configPath, err)
	}

	return nil
}
