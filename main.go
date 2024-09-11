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
	addRiotAPILimits()

	if err := limiter.LoadRequestsMade(); err != nil {
		fmt.Println("Couldn't load previous records - " + err.Error())
	}

	defer func() {
		// Should also add option for retry
		if err := limiter.SaveRequestsMade(); err != nil {
			panic("Couldn't save records. Please wait for 120 seconds before starting the program again - " + err.Error())
		}
	}()

	var user UserData
	if err := user.FetchData(); err != nil {
		panic(fmt.Sprintf("error fetching data for user %s", err.Error()))
	}

	var game GameData
	lastGameID, err := getLastGameID(user.PUUID)
	if err != nil {
		panic(fmt.Sprintf("not able to retrieve the last game IDs for PUUID %s - %s", user.PUUID, err.Error()))
	}
	if err := game.GetGameData(lastGameID); err != nil {
		panic(fmt.Sprintf("not able to retrieve the game data for game %s - %s", lastGameID, err.Error()))
	}

	gameSheet := queueFormatting(game.Info.QueueID)
	if _, err := os.Stat(user.ExcelFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating file %s with sheet %s\n", user.ExcelFile, gameSheet)
		if err = CreateGameRecordFile(user.ExcelFile, gameSheet); err != nil {
			panic(fmt.Sprintf("not able to create the excel file %s - %s", user.ExcelFile, err.Error()))
		}
	}

	fmt.Println("Adding game to sheet " + gameSheet)
	err = AddGameRecord(user.ExcelFile, gameSheet, game, user)
	if err != nil {
		panic(fmt.Sprintf("not able to add the game stats to file %s, sheet %s - %s", user.ExcelFile, gameSheet, err.Error()))
	}
}
