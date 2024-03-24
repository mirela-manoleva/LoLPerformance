package main

import (
	"fmt"
	"net/http"

	api "main/riot_api"
)

func main() {
	client := &http.Client{}
	PUUID := api.GetPUUID(client, "Alerim", "EUNE")
	lastGameID := api.GetLastGameID(client, PUUID)
	fmt.Println(lastGameID)
}
