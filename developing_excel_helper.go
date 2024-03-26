package main

/*
	Simulating a request. Freely remove this file when developing ends.
*/

/*
Mainly added for ease of use.
Much better to pass this type instead of having a function with 10-15 arguments.
*/
type GameData struct {
	unixStartTimestamp int64 //Milliseconds

	queueType int
	isWin     bool
	role      string
	champion  string

	killParticipation float64
	kills             int
	deaths            int
	assists           int
	kda               float64

	gameLengthSeconds float64
	damagePerMinute   float64
	goldPerMinute     float64
	csPerMinute       float64
}

type SummonerData struct {
	rank string
}

/*
Made in order not to make api calls every time.
Remove it when the api and JSON parser are done.
*/
func developingFunctionGetGameInfo() (GameData, SummonerData) {
	return GameData{
			unixStartTimestamp: 1725527594000,

			queueType: 420, // "Ranked"
			isWin:     false,
			role:      "ADC",
			champion:  "Jinx",

			kills:             1,
			deaths:            2,
			assists:           33,
			killParticipation: 0.6343,
			kda:               17,

			gameLengthSeconds: 1952.2575888, // 32:30
			damagePerMinute:   33.1234,
			goldPerMinute:     18,
			csPerMinute:       0,
		},
		SummonerData{
			rank: "Bronze 1",
		}
}
