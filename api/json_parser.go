package api

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
	Participants []ParticipantData `json:"participants"`
	QueueId      int               `json:"queueId"`
}

type ParticipantData struct {
	Assists                     int        `json:"assists"`
	Challenges                  Challenges `json:"challenges"`
	ChampionName                string     `json:"championName"`
	Deaths                      int        `json:"deaths"`
	GoldEarned                  int        `json:"goldEarned"`
	Kills                       int        `json:"kills"`
	Lane                        string     `json:"lane"`
	SummonerId                  string     `json:"summonerId"`
	TotalDamageDealtToChampions int        `json:"totalDamageDealtToChampions"`
	TotalMinionsKilled          int        `json:"totalMinionsKilled"`
	Win                         bool       `json:"win"`
}

type Challenges struct {
	DamagePerMinute     float32 `json:"damagePerMinute"`
	DeathsByEnemyChamps int     `json:"deathsByEnemyChamps"`
	GameLength          float32 `json:"gameLength"`
	GoldPerMinute       float32 `json:"goldPerMinute"`
	Kda                 float32 `json:"kda"`
	KillParticipation   float32 `json:"killParticipation"`
	Takedowns           int     `json:"takedowns"`
}

func parsePUUIDToObj(jsonStr string) (jsonObj PUUIDObj) {
	dec := json.NewDecoder(strings.NewReader(jsonStr))

	var respObject PUUIDObj
	err := dec.Decode(&respObject)
	if err != nil {
		fmt.Println("Error while decoding json string: ", err)
	}

	return respObject
}
func parseGameInfoToObj(jsonStr string) (jsonObj GameInfoObject) {
	var gameInfoObj GameInfoObject
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	err := dec.Decode(&gameInfoObj)
	if err != nil {
		fmt.Println("Error while decoding json string: ", err)
	}

	return gameInfoObj
}
