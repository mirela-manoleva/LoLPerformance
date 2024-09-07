package main

/*
	File description:
	Defines the functions that communicate with the riot servers to fetch the data needed.
	Uses the limiter to regulate api calls.
*/

import (
	"errors"
	"fmt"
	"main/limiter"
	"math"
	"net/http"
	"regexp"
	"time"
)

const (
	// This api key is for devs only. Change or research before you send to others
	DEV_API_KEY      = "RGAPI-5477a09f-ed9a-482f-b672-a364ce6d8015"
	RIOT_SERVER_EU   = "https://europe.api.riotgames.com"
	RIOT_SERVER_EUNE = "https://eun1.api.riotgames.com"
)

const (
	PUUID_ENDPOINT          = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	LAST_GAME_ID_ENDPOINT   = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	LAST_GAME_INFO_ENDPOINT = "/lol/match/v5/matches/%s"
	RANK_ENDPOINT           = "/lol/league/v4/entries/by-summoner/%s"
)

const (
	GAME_ID_REGEX = "[^\"\\[]+"
)

type GameRecord struct {
	date string
	rank string

	queueType string
	outcome   string
	role      string
	champion  string

	killParticipation float64
	kills             int
	deaths            int
	assists           int
	kda               float64

	gameLength      string
	damagePerMinute float64
	goldPerMinute   float64
	csPerMinute     float64
}

/*
The default client + 30 sec timeout
*/
var httpClient = &http.Client{Timeout: 30 * time.Second}

/*
Defines the riot API limits.
20 requests every 1 seconds.
100 requests every 2 minutes.
As there are multiple limitations those limits are the most strict ones, ignoring endpoint specific ones
*/
func addRiotAPILimits() {
	limiter.AddLimit(20, time.Second)
	limiter.AddLimit(100, 2*time.Minute)
}

func sendRiotAPIRequest(requestType string, url string) (string, error) {
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return "", errors.New("error in creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}
	request.Header.Add("X-Riot-Token", DEV_API_KEY)

	response, err := limiter.SendRequest(httpClient, request)
	if err != nil {
		return "", err
	}

	return response, nil
}

func GetPUUID(gameName string, tagLine string) (string, error) {
	requestType := "GET"
	url := RIOT_SERVER_EU + PUUID_ENDPOINT
	url = fmt.Sprintf(url, gameName, tagLine)

	responseStr, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	var respObj PUUIDObj
	err = parseStringToJSON(responseStr, &respObj)
	if err != nil {
		return "", err
	}

	return respObj.PUUID, nil
}

func GetLastGameID(PUUID string) (string, error) {
	requestType := "GET"
	url := RIOT_SERVER_EU + LAST_GAME_ID_ENDPOINT
	url = fmt.Sprintf(url, PUUID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	return regexp.MustCompile(GAME_ID_REGEX).FindString(response), nil
}

func GetGameRecord(gameID string, PUUID string) (GameRecord, error) {
	var gameRecord GameRecord

	gameInfoObj, err := getGameInfoObj(gameID)
	if err != nil {
		return gameRecord, err
	}

	participantIndex := getParticipantIndex(gameInfoObj.Metadata.Participants, PUUID)
	if participantIndex == -1 {
		return gameRecord, errors.New("error: couldn't find player in participants array")
	}
	player := gameInfoObj.Info.Participants[participantIndex]

	gameRecord.date = getDateStr(gameInfoObj.Info.GameStartTimestamp)
	gameRecord.queueType = getQueueStr(gameInfoObj.Info.QueueId)
	gameRecord.outcome = getGameOutcome(player.Win)
	gameRecord.role = player.IndividualPosition
	gameRecord.champion = player.ChampionName

	gameRecord.killParticipation = player.Challenges.KillParticipation
	gameRecord.kills = player.Kills
	gameRecord.deaths = player.Deaths
	gameRecord.assists = player.Assists
	gameRecord.kda = player.Challenges.Kda

	gameRecord.gameLength = getGameLength(player.Challenges.GameLength)
	gameRecord.damagePerMinute = player.Challenges.DamagePerMinute
	gameRecord.goldPerMinute = player.Challenges.GoldPerMinute
	gameRecord.csPerMinute = float64(player.TotalMinionsKilled) / player.Challenges.GameLength

	rank, err := getRank(player.SummonerId)
	if err != nil {
		return gameRecord, err
	}
	gameRecord.rank = rank

	return gameRecord, nil
}

func getGameInfoObj(gameID string) (GameInfoObject, error) {
	requestType := "GET"
	url := RIOT_SERVER_EU + LAST_GAME_INFO_ENDPOINT
	url = fmt.Sprintf(url, gameID)

	responseStr, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return GameInfoObject{}, err
	}

	var respObj GameInfoObject
	err = parseStringToJSON(responseStr, &respObj)

	return respObj, err
}

func getRank(summonerID string) (string, error) {
	requestType := "GET"
	url := RIOT_SERVER_EUNE + RANK_ENDPOINT
	url = fmt.Sprintf(url, summonerID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func getDateStr(unixTimestamp int64) string {
	seconds := unixTimestamp / 1000
	nanoseconds := (unixTimestamp % 1000) * 1000
	year, month, day := time.Unix(seconds, nanoseconds).Date()
	return fmt.Sprintf("%d/%d/%d", day, month, year)
}

func getQueueStr(queueType int) string {
	switch queueType {
	case 400:
		return "Normal"
	case 420:
		return "Ranked Solo"
	case 440:
		return "Ranked Flex"
	case 490:
		return "Normal"
	default:
		return "Other"
	}
}

func getParticipantIndex(allPUUIDs []string, playerPUUID string) int {
	for index, PUUID := range allPUUIDs {
		if PUUID == playerPUUID {
			return index
		}
	}
	return -1
}

func getGameOutcome(isWon bool) string {
	if isWon {
		return "Win"
	}
	return "Loss"
}

func getGameLength(lengthInSeconds float64) string {
	minutes := int64(lengthInSeconds / 60)
	seconds := int64(math.Round(lengthInSeconds - float64(minutes*60)))
	return fmt.Sprintf("%d:%d", minutes, seconds)
}
