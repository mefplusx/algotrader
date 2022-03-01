package bootstrap

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPAPI_PORT string `yaml:"HTTPAPI_PORT"`

	HOSTCONNECTOR_ADDRESS string `yaml:"HOSTCONNECTOR_ADDRESS"`

	EXCHANGE_API_KEY        string `yaml:"EXCHANGE_API_KEY"`
	EXCHANGE_API_SECRET     string `yaml:"EXCHANGE_API_SECRET"`
	EXCHANGE_API_PASSPHRASE string `yaml:"EXCHANGE_API_PASSPHRASE"`
}

func (c *Config) Read(fileName string) {
	path, _ := os.Getwd()
	b, _ := os.ReadFile(path + "/" + fileName)
	yaml.Unmarshal(b, c)
}
