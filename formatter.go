package main

/*
	File description:
	Defines the functionality that formats the data to prepare it for presentation in the EXCEL file.
*/

import "errors"

type PlayerGameData struct {
	unixStartTimestamp int64 //Milliseconds

	queueType int
	win       bool
	role      string
	champion  string

	killParticipation float64
	kills             int
	deaths            int
	assists           int
	kda               float64

	gameLengthSeconds float64
	damagePerMinute   float64
	goldPerMinute     float64
	csPerMinute       float64
}

func GetGameRecord(gameID string, PUUID string) (PlayerGameData, Rank, error) {
	var data PlayerGameData

	gameInfoObj, err := getGameData(gameID)
	if err != nil {
		return data, Rank{}, err
	}

	participantIndex, err := getParticipantIndex(gameInfoObj.Metadata.ParticipantsPUUID, PUUID)
	if err != nil {
		return data, Rank{}, err
	}
	player := gameInfoObj.Info.ParticipantData[participantIndex]

	data.unixStartTimestamp = gameInfoObj.Info.GameStartTimestamp

	data.queueType = gameInfoObj.Info.QueueID
	data.win = player.Win
	data.role = player.TeamPosition
	data.champion = player.Champion

	data.killParticipation = player.Challenges.KP
	data.kills = player.Kills
	data.deaths = player.Deaths
	data.assists = player.Assists
	data.kda = player.Challenges.KDA

	data.gameLengthSeconds = player.Challenges.GameLength
	data.damagePerMinute = player.Challenges.DPM
	data.goldPerMinute = player.Challenges.GPM
	data.csPerMinute = (float64(player.CS) / player.Challenges.GameLength) * 60

	rank, err := getRank(player.SummonerID)

	return data, rank, err
}

func getParticipantIndex(allPUUIDs []string, playerPUUID string) (int, error) {
	for index, PUUID := range allPUUIDs {
		if PUUID == playerPUUID {
			return index, nil
		}
	}
	return -1, errors.New("error: couldn't find player in participants array")
}
