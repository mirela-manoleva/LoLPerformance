package limiter

import (
	"net/http"
	"testing"
	"time"
)

/*
The default client + 30 sec timeout
*/
var httpClient = &http.Client{Timeout: 30 * time.Second}

func TestRegRequestPass(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	limitRequests := 10
	limitPeriod := time.Second
	AddLimit(limitRequests, limitPeriod)

	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	for i := 0; i < limitRequests-1; i++ {
		addRecord(time.Now())
	}

	response, err := SendRequest(httpClient, request)
	if err != nil || len(response) == 0 {
		t.Fatal("A request failed or the response was with len 0: " + err.Error() + "\n" + getState())
	}
}

func TestRegRequestPassRealRequests(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	limitRequests := 2
	limitPeriod := time.Second
	AddLimit(limitRequests, limitPeriod)

	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	for i := 0; i < limitRequests; i++ {
		response, err := SendRequest(httpClient, request)
		if err != nil || len(response) == 0 {
			t.Fatal("A request failed or the response was with len 0: " + err.Error() + "\n" + getState())
		}
	}
}

func TestRegRequestFailRealRequests(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	limitRequests := 3
	limitPeriod := time.Second
	AddLimit(limitRequests, limitPeriod)
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	for i := 0; i < limitRequests; i++ {
		response, err := SendRequest(httpClient, request)
		if err != nil || len(response) == 0 {
			t.Fatal("A request failed or the response was with len 0: " + err.Error() + "\n" + getState())
		}
	}

	response, err := SendRequest(httpClient, request)
	if err == nil {
		t.Fatal("A request succeeded when it shouldn't have.\n" + getState())
	}
	if len(response) != 0 {
		t.Fatal("The response length isn't 0.")
	}
}

func TestRegRequestFail(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	limitRequests := 6
	limitPeriod := time.Second
	AddLimit(limitRequests, limitPeriod)

	for i := 0; i < limitRequests; i++ {
		addRecord(time.Now())
	}

	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}
	response, err := SendRequest(httpClient, request)
	if err == nil {
		t.Fatal("A request succeeded when it shouldn't have.\n" + getState())
	}
	if len(response) != 0 {
		t.Fatal("The response length isn't 0.")
	}
}

func TestRegRequestTwoLimitsPass(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	limit1Requests := 2
	limit1Period := 10 * time.Millisecond
	AddLimit(limit1Requests, limit1Period)

	limit2Requests := 2*limit1Requests - 1
	limit2Period := 2 * limit1Period
	AddLimit(limit2Requests, limit2Period)

	for i := 0; i < limit1Requests-1; i++ {
		addRecord(time.Now())
	}

	response, err := SendRequest(httpClient, request)
	if err != nil || len(response) == 0 {
		t.Fatal("A request failed or the response was with len 0: " + err.Error() + "\n" + getState())
	}

	clearAllRecords()

	for i := 0; i < limit2Requests-1; i++ {
		addRecord(time.Now()) // Technically can break the first limit but we only care for the second one
	}

	time.Sleep(limit1Period) // Ensure that we won't break the first limit
	response, err = SendRequest(httpClient, request)
	if err != nil || len(response) == 0 {
		t.Fatal("A request failed or the response was with len 0: " + err.Error() + "\n" + getState())
	}
}

func TestRegRequestTwoLimitsFail(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	url := "https://www.google.com"
	requestType := "GET"

	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		t.Fatal("error creating a request [" + requestType + ", " + url + "] - " + err.Error())
	}

	limit1Requests := 2
	limit1Period := 10 * time.Millisecond
	AddLimit(limit1Requests, limit1Period)

	limit2Requests := 2*limit1Requests - 1
	limit2Period := 2 * limit1Period
	AddLimit(limit2Requests, limit2Period)

	for i := 0; i < limit1Requests; i++ {
		addRecord(time.Now())
	}

	response, err := SendRequest(httpClient, request)
	if err == nil {
		t.Fatal("A request succeeded when it shouldn't have.\n" + getState())
	}
	if len(response) != 0 {
		t.Fatal("The response length isn't 0.")
	}

	clearAllRecords()

	for i := 0; i < limit2Requests; i++ {
		addRecord(time.Now()) // Technically can break the first limit but we only care for the second one
	}

	time.Sleep(limit1Period) // Ensure that we won't break the first limit
	response, err = SendRequest(httpClient, request)
	if err == nil {
		t.Fatal("A request succeeded when it shouldn't have.\n" + getState())
	}
	if len(response) != 0 {
		t.Fatal("The response length isn't 0.")
	}
}
