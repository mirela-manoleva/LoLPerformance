package limiter

import (
	"errors"
	"io"
	"net/http"
	"time"
)

/*
	Pretty much the default client + a 30 sec timeout
*/
var customClient = &http.Client{Timeout: 30*time.Second}

/*
	Checks if the request will break any set limit and executes it.
*/
func SendRegulatedRequest(request *http.Request) (payload string, err error) {
	if request == nil {
		return "", errors.New("request is nil")
	}

	if !canExecuteRequestNow() {
		return "", errors.New("rate exceeded for API calls")
	}

	addRecord(time.Now())

	response, err := customClient.Do(request)
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
