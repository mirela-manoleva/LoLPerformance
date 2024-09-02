package main

import (
	"main/limiter"
	"time"
)

func main() {
	/*
		20 requests every 1 seconds
		100 requests every 2 minutes
		As there are multiple limitations those limits will use the most strict ones, ignoring endpoint specific ones
	*/
	limiter.AddLimit(20, time.Second)
	limiter.AddLimit(100, 2*time.Minute)

	err := limiter.LoadRequestsMade()
	if err != nil {
		println("Couldn't load previous records - " + err.Error())
	}

	defer func () {
			err = limiter.SaveRequestsMade()
		// Should also add option for retry
		if err != nil {
			println("Couldn't save records. Please wait for 120 seconds before starting the program again - " + err.Error())
		}
	} ()

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

	lastGameInfo, err := GetLastGameInfo(lastGameID)
	if err != nil {
		println("GetLastGameInfo Error: ", err.Error())
		return
	}

	println(lastGameInfo)
}