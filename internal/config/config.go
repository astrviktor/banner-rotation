package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPServer HTTPServerConfig
	DB         DBConfig
	Kafka      KafkaConfig
}

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	DSN                string `yaml:"dsn"`
	MaxConnectAttempts int    `yaml:"maxConnectAttempts"`
}

type KafkaConfig struct {
	Topic              string `yaml:"topic"`
	BrokerAddress      string `yaml:"brokerAddress"`
	MaxConnectAttempts int    `yaml:"maxConnectAttempts"`
}

func NewConfig(name string) Config {
	var config Config

	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err.Error())
		return DefaultConfig()
	}

	err = yaml.Unmarshal(file, &config)
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
		DBConfig{DSN: "postgres://user:password@postgres:5432/banner_rotation", MaxConnectAttempts: 5},
		KafkaConfig{Topic: "events", BrokerAddress: "kafka:9092", MaxConnectAttempts: 5},
	}
}
