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
	"net/http"
	"regexp"
	"time"
)

const (
	// This api key is for devs only. Change or research before you send to others
	DEV_API_KEY      = "RGAPI-5477a09f-ed9a-482f-b672-a364ce6d8015"
	RIOT_SERVER_EU   = "https://europe.api.riotgames.com"
	RIOT_SERVER_EUNE = "https://eun1.api.riotgames.com"
	RIOT_SERVER_EUW  = "https://euw1.api.riotgames.com"
)

const (
	PUUID_ENDPOINT         = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	LAST_GAME_ID_ENDPOINT  = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	GAME_DATA_ENDPOINT     = "/lol/match/v5/matches/%s"
	SUMMONER_DATA_ENDPOINT = "/lol/league/v4/entries/by-summoner/%s"
)

const (
	GAME_ID_REGEX = "[^\"\\[]+"
)

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

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	var puuid PUUID
	err = JSONToObject(response, &puuid)
	if err != nil {
		return "", err
	}

	if len(puuid.String) == 0 {
		return "", fmt.Errorf("didn't find player with name %s:%s", gameName, tagLine)
	}

	return puuid.String, nil
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

func getGameData(gameID string) (GameData, error) {
	requestType := "GET"
	url := RIOT_SERVER_EU + GAME_DATA_ENDPOINT
	url = fmt.Sprintf(url, gameID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return GameData{}, err
	}

	var data GameData
	err = JSONToObject(response, &data)

	return data, err
}

func getRank(summonerID string) (Rank, error) {
	requestType := "GET"
	url := RIOT_SERVER_EUNE + SUMMONER_DATA_ENDPOINT
	url = fmt.Sprintf(url, summonerID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return Rank{}, err
	}

	var ranks []Rank
	err = JSONToObject(response, &ranks)
	if err != nil {
		return Rank{}, err
	}

	for i := 0; i < len(ranks); i++ {
		if ranks[i].QueueType == "RANKED_SOLO_5x5" {
			return ranks[i], nil
		}
	}

	// if no soloqueue rank is found
	return Rank{Name: "UNRANKED", Number: ""}, nil
}
