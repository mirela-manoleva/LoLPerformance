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

func main() {
	getConfig()
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

	if _, err := os.Stat(excelFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating file %s with sheet %s\n", excelFile, gameSheet)
		if err = CreateGameRecordFile(excelFile, gameSheet); err != nil {
			panic(fmt.Sprintf("not able to create the excel file %s - %s", excelFile, err.Error()))
		}
	}

	fmt.Println("Adding game to sheet " + gameSheet)
	err = AddGameRecord(excelFile, gameSheet, gameRecord, rank, summonerName+"#"+summonerTag)
	if err != nil {
		panic(fmt.Sprintf("not able to add the game stats to file %s, sheet %s - %s", excelFile, gameSheet, err.Error()))
	}
}
