package limiter

import (
	"strconv"
	"strings"
	"time"
)

func limitsToString(slice []limit) string {
	var result strings.Builder
	for i := 0; i < len(slice); i++ {
		result.WriteString("Limit " + strconv.Itoa(i+1) + " : Req = " + strconv.Itoa(slice[i].requestCount) + " Per = " + slice[i].period.String() + " seconds." + "\n")
	}
	return result.String()
}

func recordsToString(slice []time.Time) string {
	var result strings.Builder
	for i := 0; i < len(slice); i++ {
		result.WriteString("Request " + strconv.Itoa(i+1) + " done at time: " + slice[i].String() + "\n")
	}
	return result.String()
}

func getState() string {
	var result strings.Builder
	result.WriteString("State:\n")
	result.WriteString("Limits:\n" + limitsToString(requestLimits))
	result.WriteString("Records:\n" + recordsToString(requestRecords))
	result.WriteString("Time now:" + time.Now().String() + "\n")
	return result.String()
}
