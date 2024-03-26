package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

func makeRequestRiotAPI (requestType string, url string) (*http.Request, error) {
	request, err := http.NewRequest(requestType, url, nil)
	if(err != nil) {
		return nil, errors.New("error in creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	return request, nil
}

func GetPUUID(gameName string, tagLine string) (string, error) {
	requestType := "GET"
	url := EU_HOST + API_GET_PUUID
	url = fmt.Sprintf(url, gameName, tagLine)

	request, err := makeRequestRiotAPI(requestType, url)
	if(err != nil){
		return "", err
	}

	response, err := SendRegulatedRequest(request)
	if (err != nil){
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

	request, err := makeRequestRiotAPI(requestType, url)
	if(err != nil){
		return "", err
	}

	response, err := SendRegulatedRequest(request)
	if(err != nil){
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

	request, err := makeRequestRiotAPI(requestType, url)
	if(err != nil){
		return "", err
	}

	response, err := SendRegulatedRequest(request)
	if(err != nil){
		return "", err
	}

	return response, nil
}

func GetRank(summonerID string) (string, error) {
	requestType := "GET"
	url := API_GET_RANK
	url = fmt.Sprintf(url, summonerID)

	request, err := makeRequestRiotAPI(requestType, url)
	if(err != nil){
		return "", err
	}

	response, err := SendRegulatedRequest(request)
	if(err != nil){
		return "", err
	}

	return response, nil
}
