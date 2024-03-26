/*
	Toto note:
	This file has functionallity mainly used for debugging and testing. At some point it will be reworked into tests.
	There is no reason to use any of the functions written here it in the application.
	If the app needs a function from this file just move the function to some other file.
	This is the only part of the code that has tests because it is important we don't abuse the Riot API.
*/

package api

import (
	"slices"
	"time"
)

func DebugPrintAllLimits() {
	for i := 0; i < len(requestLimits); i++ {
		println("Limit ", i+1, " : Req = ", requestLimits[i].requestsAllowed, " Per = ", requestLimits[i].period.Seconds(), " seconds.")
	}
}

func DebugPrintAllRecordedRequests() {
	for i := 0; i < len(requestRecords); i++ {
		println("Request ", i+1, " done at time: ", requestRecords[i].String())
	}
}

// Can be useful for simulating certain events.
func DebugAddRecord(record time.Time) {
	i := 0
	for ; i < len(requestRecords) ; i++ {
		if(record.Before(requestRecords[i])) {
			requestRecords = slices.Insert(requestRecords, i, record)
			return
		}
	}

	requestRecords = append(requestRecords, record)
}

func TestSingleLimit() {
	println("Starting test with single limit")

	limitResponse := 5
	limitPeriod := time.Second

	// Check without limits
	DebugAddRecord(time.Now())

	if(!canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return true. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Sends request when no limits are specified.")
	}

	// Check with empty records
	clearAllRecords()

	AddLimit(limitResponse, limitPeriod)
	if(!canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return true. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Sends request when there are no records.")
	}

	// Check with limitResponse - 1
	clearAllRecords()
	for i := 0; i < limitResponse - 1; i++ {
		DebugAddRecord(time.Now())
	}

	if(!canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return true. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Send a request that doesn't break a limit.")
	}

	// Check with limitResponse
	clearAllRecords()
	for i := 0; i < limitResponse; i++ {
		DebugAddRecord(time.Now())
	}

	if(canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return false. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Didn't send a request that would break a limit.")
	}

	// Check that canExecuteRequestNow() clears the records if enough time has passed
	time.Sleep(limitPeriod)
	canExecuteRequestNow()
	if(len(requestRecords) != 0){
		println("Error, canExecuteRequestNow() should have cleared the records. Info:")
		println("Limits after wait:")
		DebugPrintAllLimits()
		println("Records after wait:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! canExecuteRequestNow() cleared the records when it was supposed to.")
	}

	// Check that canExecuteRequestNow() doesn't clear the records when it shouldn't
	clearAllRecords()
	for i := 0; i < limitResponse; i++ {
		DebugAddRecord(time.Now())
	}
	numberOfRecords := len(requestRecords)
	time.Sleep(limitPeriod - 250*time.Millisecond)

	canExecuteRequestNow()
	if(len(requestRecords) != numberOfRecords){
		println("Error, canExecuteRequestNow() shouldn't have cleared any records. Info:")
		println("Limits after wait:")
		DebugPrintAllLimits()
		println("Records after wait:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! canExecuteRequestNow() didn't clear the records when it wasn't supposed to.")
	}
}

func TestSaveLoadRecords() {
	println("Starting tests for Save and Load records.")

	numberOfRequest := 5
	timeBetweenRequests := 50*time.Millisecond // just to make is easier to read

	clearAllRecords()
	for i := 0; i < numberOfRequest; i++ {
		time.Sleep(timeBetweenRequests)
		DebugAddRecord(time.Now())
	}

	tmp := make([]time.Time, len(requestRecords))
	copy(tmp, requestRecords)

	err := SaveRecords()
	if err != nil {
		println("Error, SaveRecords returned error - " + err.Error())
		return
	}

	clearAllRecords()

	err = LoadRecords()
	if err != nil {
		println("Error, LoadRecords returned error - " + err.Error())
		return
	}

	areDifferent := false

	if len(tmp) != len(requestRecords) {
		areDifferent = true
		return
	}

	for i := 0; i < len(requestRecords); i++ {
		if(!tmp[i].Equal(requestRecords[i])) {
			areDifferent = true
			return
		}
	}

	defer func () {
			if(areDifferent) {
			println("Error, there are differences before and after the save-load functions. Info:")
			println("Records before the operations:")
			for i := 0; i < len(tmp); i++ {
				println("Request ", i+1, " done at time: ", tmp[i].String())
			}
			println("Records after the operations:")
			DebugPrintAllRecordedRequests()
			return
		}else{
			println("Success! The records before and after the save-load functions are the same.")
		}
	}()
}

func TestSaveLoadWithLimit() {
	println("Starting tests for Save and Load and a set limit.")

	numberOfRequests := 5
	timeBetweenRequests := 50*time.Millisecond // just to make is easier to read

	AddLimit(numberOfRequests + 1, time.Second)

	clearAllRecords()
	for i := 0; i < numberOfRequests; i++ {
		time.Sleep(timeBetweenRequests)
		DebugAddRecord(time.Now())
	}

	err := SaveRecords()
	if err != nil {
		println("Error, SaveRecords returned error - " + err.Error())
		return
	}

	clearAllRecords()

	err = LoadRecords()
	if err != nil {
		println("Error, LoadRecords returned error - " + err.Error())
		return
	}

	// It should be able to add one more request
	if(!canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return true. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Send a request that doesn't break a limit.")
	}

	err = SaveRecords()
	if err != nil {
		println("Error, SaveRecords returned error - " + err.Error())
		return
	}

	clearAllRecords()

	err = LoadRecords()
	if err != nil {
		println("Error, LoadRecords returned error - " + err.Error())
		return
	}

	DebugAddRecord(time.Now())

	// Should not be able to add another
	if(canExecuteRequestNow()){
		println("Error, canExecuteRequestNow() should return false. Info:")
		println("Limits:")
		DebugPrintAllLimits()
		println("Records:")
		DebugPrintAllRecordedRequests()
		println("Current time: ", time.Now().String())
		return
	} else {
		println("Success! Didn't send a request that would break a limit.")
	}
}