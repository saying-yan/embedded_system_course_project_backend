package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type WebConf struct {
	Port     int    `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type ConnectorConf struct {
	Port            int `yaml:"port"`
	TimeoutDuration int `yaml:"timeout_duration"`
}

type LoggerConf struct {
	Path   string `yaml:"path"`
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

type Config struct {
	WebConf       `yaml:"web"`
	ConnectorConf `yaml:"connector"`
	LoggerConf    `yaml:"logger"`
}

func LoadConfig(configFile string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	if config.WebConf.Port == 0 {
		config.WebConf.Port = 80
	}
	if config.ConnectorConf.Port == 0 {
		config.ConnectorConf.Port = 8888
	}

	return config, nil
}
