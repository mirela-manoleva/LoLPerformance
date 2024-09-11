package main

/*
	File description:
	Helper functions that are used for formatting the values from the Riot client.
	Most of them are API specific so it makes sense to be in a riot specific file.
*/

/*
Riot API returns a number as a queue type.
This function returns the queue type in a a normal format.
Reference on them can be found here:
https://static.developer.riotgames.com/docs/lol/queues.json
*/
func queueFormatting(queueType int) string {
	switch queueType {
	case 400:
		return "Draft Pick"
	case 420:
		return "Ranked Solo"
	case 430:
		return "Blind Pick"
	case 440:
		return "Ranked Flex"
	case 450:
		return "ARAM"
	case 490:
		return "Quickplay"
	case 700:
		return "Clash"
	case 720:
		return "ARAM Clash"
	default:
		return "Other"
	}
}

/*
Riot API returns bool type.
This function returns the outcome in a string format.
*/
func outcomeFormatting(isWin bool) string {
	if isWin {
		return "Win"
	} else {
		return "Loss"
	}
}

/*
Riot API returns weird string.
This function returns the role in a normal string format.
*/
func roleFormatting(role string) string {
	if role == "TOP" {
		return "Top"
	} else if role == "JUNGLE" {
		return "Jungle"
	} else if role == "MIDDLE" {
		return "Mid"
	} else if role == "BOTTOM" {
		return "ADC"
	} else {
		return "Support"
	}
}
