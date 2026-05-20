package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	EmailPassword = "email-password"
	WhiteListCode = "whitelist-code"
	KafkaBrokers  = "kafka-brokers"
)

var conf *Config

type Config struct {
	EmailPassword string
	WhitelistCode string `yaml:"whitelist-code"`
	KafkaBrokers  string `yaml:"kafka-brokers"`
}

func Init() error {
	ENV := os.Getenv("ENV")

	body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	if err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	conf.EmailPassword = os.Getenv("EMAIL_PASSWORD")

	return nil
}

func Get(key string) interface{} {
	switch key {
	case EmailPassword:
		return conf.EmailPassword
	case WhiteListCode:
		return conf.WhitelistCode
	case KafkaBrokers:
		return strings.Split(conf.KafkaBrokers, ",")
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
