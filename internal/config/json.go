package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFileName string = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%s/%s", homeDir, configFileName)
	return path
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func writeConfig(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(getConfigFilePath(), data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	err := writeConfig(*c)
	if err != nil {
		return err
	}

	return nil
}
