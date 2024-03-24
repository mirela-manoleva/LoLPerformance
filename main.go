package main

import (
	"fmt"
	"main/api"
	"net/http"
)

func main() {
	client := &http.Client{}
	PUUID := api.GetPUUID(client, "Alerim", "EUNE")
	lastGameID := api.GetLastGameID(client, PUUID)
	lastGameInfo := api.GetLastGameInfo(client, lastGameID)
	AlerimRank := api.GetRank(client, "LeLTvLbZnC6-6NyvEj9aC2hCUJV3O1iThzV9YQqVtDkHe7E")
	fmt.Println(lastGameInfo)
	fmt.Println(AlerimRank)
}
