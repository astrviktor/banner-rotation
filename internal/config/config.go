package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPServer HTTPServerConfig
}

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var GlobalConfig Config

func NewConfig(file string) Config {
	var config Config

	yamlFile, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return DefaultConfig()
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
		return DefaultConfig()
	}

	return config
}

func DefaultConfig() Config {
	log.Println("get default config")

	return Config{
		HTTPServerConfig{Host: "", Port: "8888"},
	}
}
