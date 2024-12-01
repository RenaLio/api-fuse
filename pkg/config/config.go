package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ProviderData struct {
	Name    string   `yaml:"name"`
	URL     string   `yaml:"url"`
	Models  []string `yaml:"models"`
	APIKeys []string `yaml:"api_keys"`
}

type Provider struct {
	Name string         `yaml:"provider"`
	Data []ProviderData `yaml:"data"`
}

type Config struct {
	Syntax       float64             `yaml:"syntax"`
	CheckModels  []string            `yaml:"check_models"`
	ModelMapping map[string][]string `yaml:"model_mapping"`
	Providers    []Provider          `yaml:"providers"`
}

func Load(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
