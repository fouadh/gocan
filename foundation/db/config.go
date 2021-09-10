package db

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/user"
)

type Config struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	User             string `json:"user"`
	Password         string `json:"password"`
	Database         string `json:"database"`
	Embedded         bool   `json:"embedded"`
	EmbeddedDataPath string `json:"embeddedDataPath"`
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
}

var DefaultConfig = Config{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "postgres",
	Database: "postgres",
	Embedded: true,
}

func ReadConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get user home folder")
	}
	configFile := fmt.Sprintf("%s/.gocan/config.json", usr.HomeDir)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &DefaultConfig, nil
	}
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal configuration file")
	}
	return config, nil
}
