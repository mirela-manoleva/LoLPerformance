package main

// TODO:

// Check naming. -> new branch prob
// Check the limiter again.
// Check if we need any type of concurency and if the program is concurency-safe.
// Fix error messages.

import (
	"errors"
	"fmt"
	"main/limiter"
	"os"
)

var summonerName = ""
var summonerTag = ""
var gameFile = "Improvement.xlsx"
var gameSheet = ""

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

	if _, err := os.Stat(gameFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating file %s\n", gameFile)
		if err = CreateGameRecordFile(gameFile, gameSheet); err != nil {
			panic(err)
		}
	}

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

	fmt.Println("Adding game")
	err = AddGameRecord(gameFile, gameSheet, gameRecord, rank)
	if err != nil {
		panic(err)
	}
}
