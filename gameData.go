package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

/*
	File description:
	Defines the functionality to parse the JSON responses coming from riotclient.go to Go structures
*/

const (
	LAST_GAME_ID_ENDPOINT = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	GAME_DATA_ENDPOINT    = "/lol/match/v5/matches/"
	GAME_ID_REGEX         = "[^\"\\[]+"
)

/*
Struct that holds the data about a specific LoL game.
Used when parsing the response from GAME_DATA_ENDPOINT.
*/
type GameData struct {
	Metadata struct {
		ParticipantsPUUID []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		GameStartTimestamp int64 `json:"gameStartTimestamp"`
		ParticipantData    []struct {
			Assists    int `json:"assists"`
			Challenges struct {
				DPM        float64 `json:"damagePerMinute"`
				GameLength float64 `json:"gameLength"`
				GPM        float64 `json:"goldPerMinute"`
				KDA        float64 `json:"kda"`
				KP         float64 `json:"killParticipation"`
			} `json:"challenges"`
			Champion     string `json:"championName"`
			Deaths       int    `json:"deaths"`
			Kills        int    `json:"kills"`
			TeamPosition string `json:"teamPosition"`
			SummonerID   string `json:"summonerId"`
			CS           int    `json:"totalMinionsKilled"`
			Win          bool   `json:"win"`
		} `json:"participants"`
		QueueID int `json:"queueId"`
	} `json:"info"`
}

func (game *GameData) GetGameData(gameID string) error {
	requestType := "GET"
	url := RIOT_SERVER_EU + GAME_DATA_ENDPOINT + gameID

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(strings.NewReader(response))
	err = dec.Decode(game)
	if err != nil {
		return err
	}

	if len(game.Metadata.ParticipantsPUUID) == 0 {
		return fmt.Errorf("couldn't get game [%s] data", gameID)
	}

	return nil
}

func getLastGameID(PUUID string) (string, error) {
	requestType := "GET"
	url := RIOT_SERVER_EU + LAST_GAME_ID_ENDPOINT
	url = fmt.Sprintf(url, PUUID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	lastGameID := regexp.MustCompile(GAME_ID_REGEX).FindString(response)
	if len(lastGameID) == 0 {
		return "", fmt.Errorf("couldn't get last game ID with PUUID %s", PUUID)
	}
	return lastGameID, nil
}
