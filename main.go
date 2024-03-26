package main

// TODO:

// Check naming. -> new branch prob
// Check the limiter again.
// Check if we need any type of concurency and if the program is concurency-safe.
// Fix error messages.

import (
	"main/limiter"
)

var gameFile = "Improvement.xlsx"
var gameSheet = "Mirela"

func main() {
	addRiotAPILimits()

	err := limiter.LoadRequestsMade()
	if err != nil {
		println("Couldn't load previous records - " + err.Error())
	}

	defer func() {
		err = limiter.SaveRequestsMade()
		// Should also add option for retry
		if err != nil {
			println("Couldn't save records. Please wait for 120 seconds before starting the program again - " + err.Error())
		}
	}()

	PUUID, err := GetPUUID("Alerim", "EUNE")
	if err != nil {
		println("GetPUUID Error: ", err.Error())
		return
	}

	lastGameID, err := GetLastGameID(PUUID)
	if err != nil {
		println("GetLastGameID Error: ", err.Error())
		return
	}

	gameRecord, rank, err := GetGameRecord(lastGameID, PUUID)
	if err != nil {
		println("GetLastGameInfo Error: ", err.Error())
		return
	}

	println(gameRecord.champion)
	println(rank.Name, rank.Tier)
}
