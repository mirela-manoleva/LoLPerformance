package main

import (
	"fmt"
	api "main/api"
	"time"
)
func mainLoop() {
	/*
	20 requests every 1 seconds
	100 requests every 2 minutes
	As there are multiple limitations those limits will use the most strict ones, ignoring endpoint specific ones
	*/
	api.AddLimit(20, time.Second)
	api.AddLimit(100, 2*time.Minute)

	err := api.LoadRecords()
	if(err != nil) {
		println("Couldn't load previous records (probably they didn't exist).")
	}

	PUUID, err := api.GetPUUID("Alerim", "EUNE")
	if(err != nil) {
		println("GetPUUID Error: ", err.Error())
		return
	}

	lastGameID, err := api.GetLastGameID(PUUID)
	if(err != nil) {
		println("GetLastGameID Error: ", err.Error())
		return
	}

	lastGameInfo, err := api.GetLastGameInfo(lastGameID)
	if(err != nil) {
		println("GetLastGameInfo Error: ", err.Error())
		return
	}

	AlerimRank, err := api.GetRank("LeLTvLbZnC6-6NyvEj9aC2hCUJV3O1iThzV9YQqVtDkHe7E")
	if(err != nil) {
		println("GetRank Error: ", err.Error())
		return
	}

	fmt.Println(lastGameInfo)
	fmt.Println(AlerimRank)

	err = api.SaveRecords()
	// Should also add option for retry prob
	if(err != nil) {
		println("Couldn't save records. Please wait for 120 seconds before starting the program again!")
	}
}

// The only reason to have mainLoop and main is to run tests easier
func main() {
	//mainLoop()
	api.TestSingleLimit()
	api.TestSaveLoadRecords()
	api.TestSaveLoadWithLimit()
}
