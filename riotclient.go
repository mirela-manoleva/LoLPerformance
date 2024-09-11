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
	"time"
)

const (
	// This api key is for devs only. Change or research before you send to others
	DEV_API_KEY      = "RGAPI-5477a09f-ed9a-482f-b672-a364ce6d8015"
	RIOT_SERVER_EU   = "https://europe.api.riotgames.com"
	RIOT_SERVER_EUNE = "https://eun1.api.riotgames.com"
	RIOT_SERVER_EUW  = "https://euw1.api.riotgames.com"
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

/*
Creates an HTTP request to a riot server.
Adds the API key information in the header.
*/
func sendRiotAPIRequest(requestType string, url string) (string, error) {
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return "", fmt.Errorf("error in creating a request [%s, %s] - %s", requestType, url, err.Error())
	}
	request.Header.Add("X-Riot-Token", DEV_API_KEY)

	response, err := limiter.SendRequest(httpClient, request)
	if err != nil {
		return "", err
	}

	return response, nil
}
