package main

/*
	File description:
	Defines a class that handles all the out-of-game user data.
	PUUID, SummonedID and Rank
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var userConfigFile = filepath.Join("config", "user_config.json")

type UserData struct {
	// From userConfigFile
	SummonerName string `json:"summonerName,omitempty"`
	SummonerTag  string `json:"summonerTag,omitempty"`
	Region       string `json:"region,omitempty"`
	ExcelFile    string `json:"excelFile,omitempty"`

	// From PUUID_ENDPOINT
	PUUID string `json:"puuid,omitempty"`

	// From SUMMONER_DATA_ENDPOINT
	SummonerID string `json:"id,omitempty"`

	// From RANKED_DATA_ENDPOINT
	Tier string `json:"tier,omitempty"`
	Rank string `json:"rank,omitempty"`
}

func (user *UserData) FetchData() error {
	err := user.loadConfig()
	if err != nil {
		return err
	}

	err = user.getPUUID()
	if err != nil {
		return err
	}

	err = user.getSummonerID()
	if err != nil {
		return err
	}

	err = user.getRank()
	if err != nil {
		return err
	}

	return nil
}

func (user *UserData) loadConfig() error {
	file, err := os.ReadFile(userConfigFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, user)
}

func (user *UserData) getPUUID() error {
	requestType := "GET"
	url := RIOT_SERVER_EU + PUUID_ENDPOINT
	url = fmt.Sprintf(url, user.SummonerName, user.SummonerTag)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(strings.NewReader(response))
	err = dec.Decode(user)
	if err != nil {
		return err
	}

	if len(user.PUUID) == 0 {
		return fmt.Errorf("didn't find player's PUUID with name %s#%s in EU servers", user.SummonerName, user.SummonerTag)
	}

	return nil
}

// Function should be executed after the PUUID has been fetched
func (user *UserData) getSummonerID() error {
	if len(user.PUUID) == 0 {
		return fmt.Errorf("cannot get summoner ID for user %s#%s because there is no PUUID", user.SummonerName, user.SummonerTag)
	}

	requestType := "GET"
	var riot_server string
	if user.Region == "EUNE" {
		riot_server = RIOT_SERVER_EUNE
	} else if user.Region == "EUW" {
		riot_server = RIOT_SERVER_EUW
	} else {
		return fmt.Errorf("region %s is not EUW or EUNE", user.Region)
	}
	url := riot_server + SUMMONER_DATA_ENDPOINT
	url = fmt.Sprintf(url, user.PUUID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(strings.NewReader(response))
	err = dec.Decode(user)
	if err != nil {
		return err
	}

	if len(user.SummonerID) == 0 {
		return fmt.Errorf("cannot get summoner ID for user with PUUID %s", user.PUUID)
	}

	return nil
}

// Function should be executed after the Summoner ID has been fetched
func (user *UserData) getRank() error {
	if len(user.SummonerID) == 0 {
		return fmt.Errorf("cannot get summoner ID for user %s#%s because there is no summoner ID", user.SummonerName, user.SummonerTag)
	}

	requestType := "GET"

	var riot_server string
	if user.Region == "EUNE" {
		riot_server = RIOT_SERVER_EUNE
	} else if user.Region == "EUW" {
		riot_server = RIOT_SERVER_EUW
	} else {
		return fmt.Errorf("region %s is not EUW or EUNE", user.Region)
	}

	url := riot_server + RANKED_DATA_ENDPOINT
	url = fmt.Sprintf(url, user.SummonerID)

	response, err := sendRiotAPIRequest(requestType, url)
	if err != nil {
		return err
	}

	/*
		The request returns an array of all ranked queues that the user has played.
		We only care about the SoloQ rank.
	*/
	type RankJSON struct {
		QueueType string `json:"queueType"`
		Tier      string `json:"tier,omitempty"`
		Rank      string `json:"rank,omitempty"`
	}

	var ranks []RankJSON
	dec := json.NewDecoder(strings.NewReader(response))
	err = dec.Decode(&ranks)
	if err != nil {
		return err
	}

	for i := 0; i < len(ranks); i++ {
		if ranks[i].QueueType == "RANKED_SOLO_5x5" {
			user.Tier = ranks[i].Tier
			user.Rank = ranks[i].Rank
			return nil
		}
	}

	// if no soloqueue rank is found
	user.Tier = "UNRANKED"
	user.Rank = ""
	return nil
}
