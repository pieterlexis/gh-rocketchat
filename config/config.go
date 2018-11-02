package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func NewConfig(fname string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return parseConfig(yamlFile)
}

// parseConfig takes a yaml file as []byte and converts it to a Config
func parseConfig(configBytes []byte) (*Config, error) {
	var c Config

	c.ListenAddress = "127.0.0.1:8000"

	err := yaml.Unmarshal(configBytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
