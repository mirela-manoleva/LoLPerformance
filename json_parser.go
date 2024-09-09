package main

/*
	File description:
	Defines the functionality to parse the JSON responses coming from riotclient.go to Go structures
*/

import (
	"encoding/json"
	"strings"
)

/*
Struct that holds the riot global identification number for a player.
Used when parsing the response from PUUID_ENDPOINT.
*/
type PUUID struct {
	String string `json:"puuid"`
}

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

/*
Struct that holds the data about player's rank.
Used when parsing the response from SUMMONER_DATA_ENDPOINT.
*/
type Rank struct {
	QueueType string `json:"queueType"`
	Number    string `json:"rank,omitempty"`
	Name      string `json:"tier,omitempty"`
}

func JSONToObject(jsonStr string, object any) error {
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	err := dec.Decode(&object)
	if err != nil {
		return err
	}

	return nil
}
