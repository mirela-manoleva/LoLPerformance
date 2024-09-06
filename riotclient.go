package main

/*
	File description:
	Defines the functions that communicate with the riot servers to fetch the data needed.
	Uses the limiter to regulate api calls.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/limiter"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// This api key is for devs only. Change or research before you send to others
	TOOL_API_KEY           = "RGAPI-5477a09f-ed9a-482f-b672-a364ce6d8015"
	EU_HOST                = "https://europe.api.riotgames.com"
	API_GET_PUUID          = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	API_GET_LAST_GAME_ID   = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	GAME_ID_REGEX          = "[^\"\\[]+"
	API_GET_LAST_GAME_INFO = "/lol/match/v5/matches/%s"
	API_GET_RANK           = "https://eun1.api.riotgames.com/lol/league/v4/entries/by-summoner/%s"
)

/*
	The default client + 30 sec timeout
*/
var httpClient = &http.Client{Timeout: 30*time.Second}

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
	request.Header.Add("X-Riot-Token", TOOL_API_KEY)

	response, err := limiter.SendRequest(httpClient, request)
	if err != nil {
		return "", err
	}

	return response, nil
}

func GetPUUID(gameName string, tagLine string) (string, error) {
	requestType := "GET"
	url := EU_HOST + API_GET_PUUID
	url = fmt.Sprintf(url, gameName, tagLine)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	dec := json.NewDecoder(strings.NewReader(response))

	type ResponseObject struct {
		PUUID string `json:"puuid"`
	}
	var respObject ResponseObject
	err = dec.Decode(&respObject)
	if err != nil {
		return "", errors.New("error while decoding json string: " + err.Error())
	}

	return respObject.PUUID, nil
}

func GetLastGameID(PUUID string) (string, error) {
	requestType := "GET"
	url := EU_HOST + API_GET_LAST_GAME_ID
	url = fmt.Sprintf(url, PUUID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	return regexp.MustCompile(GAME_ID_REGEX).FindString(response), nil
}

/*
	Useful:

- Date & time
- Champion
- Rank
- Outcome
- DMG
- Game duration
- CS & Gold
- KDA

*/
func GetLastGameInfo(gameID string) (string, error) {
	requestType := "GET"
	url := EU_HOST + API_GET_LAST_GAME_INFO
	url = fmt.Sprintf(url, gameID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func GetRank(summonerID string) (string, error) {
	requestType := "GET"
	url := API_GET_RANK
	url = fmt.Sprintf(url, summonerID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return "", err
	}

	return response, nil
}
