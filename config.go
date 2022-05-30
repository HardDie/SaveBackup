package main

import (
	"encoding/json"
	"os"
)

var (
	config Config
)

type Config struct {
	Track           string `json:"track"`
	BackupDirectory string `json:"backupDirectory"`
	Name            string `json:"name"`
}

func InitConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return err
	}
	return nil
}
func GetConfig() Config {
	return config
}
