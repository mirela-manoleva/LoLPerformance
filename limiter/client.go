package limiter

/*
	File description:
	Defines the functionality to send and check http requests that abide the user defined limits.
*/

import (
	"errors"
	"io"
	"net/http"
	"time"
)

/*
	Checks if the request will break any set limit and executes it.
*/
func SendRequest(client *http.Client, request *http.Request) (payload string, err error) {
	if request == nil {
		return "", errors.New("request is nil")
	}

	if !canExecuteRequestNow() {
		return "", errors.New("rate exceeded for API calls")
	}

	addRecord(time.Now())

	// The time of execution of the request is greater than the time of the check.
	// That guarantees that the request can be executed if the original check passed.
	//
	// Note that this will change if we introduce any sort of concurency.
	// A request may happen between the check and the execution.
	// A mutex will have to be introduced.
	response, err := client.Do(request)
	if err != nil {
		return "", errors.New("error while executing request - " + err.Error())
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error while reading response body - " + err.Error())
	}

	return string(body), nil
}

/*
	Helper function for SendRequest
	Returns false if a request executed now will break any limit in requestLimits.
*/
func canExecuteRequestNow() bool {
	if len(requestLimits) == 0 {
		return true;
	}

	currentTime := time.Now()
	counter := 0

	for i, j := len(requestRecords) - 1, 0; i >= 0; i-- {
		counter++

		if counter >= requestLimits[j].requestCount {
			return false
		}

		elapsed := currentTime.Sub(requestRecords[i])
		if elapsed >= requestLimits[j].period {
			j++
			if j >= len(requestLimits) { // The difference between the last request and the current time is bigger than any limit's period
				clearAllRecords() // Only safe time to delete all records
				return true
			}
		}
	}

	return true
}