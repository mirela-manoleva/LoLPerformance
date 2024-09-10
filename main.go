package main

// TODO:

// Summoner name + tag + game file name, region and xslx file location should be settable in an external file.
// Make a log file and don't panic the program.
// Update the README.
// move the excel file open outside of the functions

import (
	"errors"
	"fmt"
	"main/limiter"
	"os"
)

var summonerName = "Alerim"
var summonerTag = "EUNE"
var gameFile = "MirelaGameStats.xlsx"

func main() {
	addRiotAPILimits()

	err := limiter.LoadRequestsMade()
	if err != nil {
		fmt.Println("Couldn't load previous records - " + err.Error())
	}

	defer func() {
		err = limiter.SaveRequestsMade()
		// Should also add option for retry
		if err != nil {
			panic("Couldn't save records. Please wait for 120 seconds before starting the program again - " + err.Error())
		}
	}()

	PUUID, err := GetPUUID(summonerName, summonerTag)
	if err != nil {
		panic(fmt.Sprintf("not able to retrieve the puuid - %s", err.Error()))
	}

	lastGameID, err := GetLastGameID(PUUID)
	if err != nil {
		panic(fmt.Sprintf("not able to retrieve the last game IDs for PUUID %s - %s", PUUID, err.Error()))
	}

	gameRecord, rank, err := GetGameRecord(lastGameID, PUUID)
	if err != nil {
		panic(fmt.Sprintf("not able to retrieve the game record for game %s and PUUID %s - %s", lastGameID, PUUID, err.Error()))
	}

	gameSheet := queueFormatting(gameRecord.queueType)

	if _, err := os.Stat(gameFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating file %s with sheet %s\n", gameFile, gameSheet)
		if err = CreateGameRecordFile(gameFile, gameSheet); err != nil {
			panic(fmt.Sprintf("not able to create the excel file %s - %s", gameFile, err.Error()))
		}
	}

	fmt.Println("Adding game to sheet " + gameSheet)
	err = AddGameRecord(gameFile, gameSheet, gameRecord, rank, summonerName+"#"+summonerTag)
	if err != nil {
		panic(fmt.Sprintf("not able to add the game stats to file %s, sheet %s - %s", gameFile, gameSheet, err.Error()))
	}
}
