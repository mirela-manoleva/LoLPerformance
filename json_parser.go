package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PUUIDObj struct {
	PUUID string `json:"puuid"`
}

type GameInfoObject struct {
	Metadata Metadata `json:"metadata"`
	Info     Info     `json:"info"`
}

type Metadata struct {
	Participants []string `json:"participants"`
}

type Info struct {
	GameStartTimestamp int64             `json:"gameStartTimestamp"`
	Participants       []ParticipantData `json:"participants"`
	QueueId            int               `json:"queueId"`
}

type ParticipantData struct {
	Assists            int        `json:"assists"`
	Challenges         Challenges `json:"challenges"`
	ChampionName       string     `json:"championName"`
	Deaths             int        `json:"deaths"`
	Kills              int        `json:"kills"`
	IndividualPosition string     `json:"individualPosition"`
	SummonerId         string     `json:"summonerId"`
	TotalMinionsKilled int        `json:"totalMinionsKilled"`
	Win                bool       `json:"win"`
}

type Challenges struct {
	DamagePerMinute   float64 `json:"damagePerMinute"`
	GameLength        float64 `json:"gameLength"`
	GoldPerMinute     float64 `json:"goldPerMinute"`
	Kda               float64 `json:"kda"`
	KillParticipation float64 `json:"killParticipation"`
}

func parseStringToJSON(jsonStr string, object any) error {
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	err := dec.Decode(&object)
	if err != nil {
		return fmt.Errorf("error while decoding json string: %s", err.Error())
	}

	return nil
}
