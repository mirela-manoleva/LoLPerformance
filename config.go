package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var configFile = filepath.Join("config", "user_config.json")

var summonerName string
var summonerTag string
var excelFile string

func getConfig() error {
	type ConfigJSON struct {
		Region       string `json:"region"`
		ExcelFile    string `json:"excelFile"`
		SummonerName string `json:"summonerName"`
		SummonerTag  string `json:"summonerTag"`
	}
	var config ConfigJSON

	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	summonerName = config.SummonerName
	summonerTag = config.SummonerTag
	excelFile = config.ExcelFile
	return nil
}
