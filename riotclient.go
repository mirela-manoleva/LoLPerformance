package main

/*
	File description:
	Defines functions that communicate with the riot servers to fetch the data needed.
	Uses the limiter to regulate api calls.
*/

import (
	"fmt"
	"main/limiter"
	"net/http"
	"os"
	"time"
)

const apiKeyFile = "config/api.key"

var RIOT_API_KEY string

const (
	RIOT_SERVER_EU   = "https://europe.api.riotgames.com"
	RIOT_SERVER_EUNE = "https://eun1.api.riotgames.com"
	RIOT_SERVER_EUW  = "https://euw1.api.riotgames.com"

	LAST_GAME_ID_ENDPOINT  = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	GAME_DATA_ENDPOINT     = "/lol/match/v5/matches/%s"
	PUUID_ENDPOINT         = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	SUMMONER_DATA_ENDPOINT = "/lol/summoner/v4/summoners/by-puuid/%s"
	RANKED_DATA_ENDPOINT   = "/lol/league/v4/entries/by-summoner/%s"
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

func loadAPIKey() error {
	content, err := os.ReadFile(apiKeyFile)
	if err != nil {
		return err
	}

	RIOT_API_KEY = string(content)
	return nil
}

/*
Creates an HTTP request to a riot server.
Adds the API key information in the header.
*/
func sendRiotAPIRequest(requestType string, url string) (string, error) {
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return "", fmt.Errorf("error in creating a request [%s, %s] - %s", requestType, url, err.Error())
	}
	request.Header.Add("X-Riot-Token", RIOT_API_KEY)

	response, err := limiter.SendRequest(httpClient, request)
	if err != nil {
		return "", err
	}

	return response, nil
}
