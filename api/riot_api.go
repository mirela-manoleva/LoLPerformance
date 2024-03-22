package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const (
	TOOL_API_KEY           = "" // Add key here
	EU_HOST                = "https://europe.api.riotgames.com"
	API_GET_PUUID          = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	API_GET_LAST_GAME_ID   = "/lol/match/v5/matches/by-puuid/%s/ids?count=1"
	GAME_ID_REGEX          = "[^\"\\[]+"
	API_GET_LAST_GAME_INFO = "/lol/match/v5/matches/%s"
	API_GET_RANK           = "https://eun1.api.riotgames.com/lol/league/v4/entries/by-summoner/%s"
)

func sendHTTPRequest(client *http.Client, typeReq string, url string) (payload string) {
	req, err := http.NewRequest(typeReq, url, nil)
	if err != nil {
		fmt.Println("Error while creating request: ", err)
		return ""
	}
	req.Header.Add("X-Riot-Token", TOOL_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while sending request: ", err)
		return ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body: ", err)
		return ""
	}

	return string(body)
}

func GetPUUID(client *http.Client, gameName string, tagLine string) (PUUID string) {
	typeReq := "GET"
	url := EU_HOST + API_GET_PUUID
	url = fmt.Sprintf(url, gameName, tagLine)
	response := sendHTTPRequest(client, typeReq, url)
	dec := json.NewDecoder(strings.NewReader(response))

	type ResponseObject struct {
		PUUID string `json:"puuid"`
	}
	var respObject ResponseObject
	err := dec.Decode(&respObject)
	if err != nil {
		fmt.Println("Error while decoding json string: ", err)
	}
	PUUID = respObject.PUUID

	return PUUID
}

func GetLastGameID(client *http.Client, PUUID string) (gameID string) {
	url := EU_HOST + API_GET_LAST_GAME_ID
	url = fmt.Sprintf(url, PUUID)
	response := sendHTTPRequest(client, "GET", url)

	gameID = regexp.MustCompile(GAME_ID_REGEX).FindString(response)

	return gameID
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
func GetLastGameInfo(client *http.Client, gameID string) (info string) {
	url := EU_HOST + API_GET_LAST_GAME_INFO
	url = fmt.Sprintf(url, gameID)
	response := sendHTTPRequest(client, "GET", url)

	return response
}

func GetRank(client *http.Client, summonerID string) (rank string) {
	url := API_GET_RANK
	url = fmt.Sprintf(url, summonerID)
	response := sendHTTPRequest(client, "GET", url)

	return response
}
