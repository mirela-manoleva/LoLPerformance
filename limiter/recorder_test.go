package limiter

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestSaveLoadRequestsMade(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	numberOfRequest := 5

	for i := 0; i < numberOfRequest; i++ {
		addRecord(time.Now())
	}

	copyRecords := make([]time.Time, len(requestRecords))
	copy(copyRecords, requestRecords)

	err := SaveRequestsMade()
	if err != nil {
		t.Fatal("SaveRequestsMade returned error - " + err.Error())
	}

	clearAllRecords()

	err = LoadRequestsMade()
	if err != nil {
		t.Fatal("LoadRequestsMade returned error - " + err.Error())
	}

	areDifferent := false
	defer func () {
		if areDifferent {
			t.Log("Differences before and after the save-load functions.\n")
			t.Log("Records before the operations:\n" + recordsToString(copyRecords))
			t.Log("Records after the operations:\n" + recordsToString(requestRecords))
			t.FailNow()
		}
	}()

	if len(copyRecords) != len(requestRecords) {
		areDifferent = true
		return
	}

	for i := 0; i < len(requestRecords); i++ {
		if !copyRecords[i].Equal(requestRecords[i]) {
			areDifferent = true
			return
		}
	}
}

func TestSaveLoadCanExec(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	numberOfRequests := 5

	AddLimit(numberOfRequests, time.Second)

	for i := 0; i < numberOfRequests - 1; i++ {
		addRecord(time.Now())
	}

	err := SaveRequestsMade()
	if err != nil {
		t.Fatal("SaveRequestsMade returned error - " + err.Error())
	}

	clearAllRecords()

	err = LoadRequestsMade()
	if err != nil {
		t.Fatal("LoadRequestsMade returned error - " + err.Error())
	}

	if !canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return true.\n" + getState())
	}
}

func TestSaveLoadCannotExec(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	numberOfRequests := 5

	AddLimit(numberOfRequests, time.Second)

	for i := 0; i < numberOfRequests - 1; i++ {
		addRecord(time.Now())
	}

	err := SaveRequestsMade()
	if err != nil {
		t.Fatal("SaveRequestsMade returned error - " + err.Error())
	}

	clearAllRecords()

	err = LoadRequestsMade()
	if err != nil {
		t.Fatal("LoadRequestsMade returned error - " + err.Error())
		return
	}

	addRecord(time.Now())

	if canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return false.\n" + getState())
	}
}

func TestSaveEmptyRecords(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	if err := os.Remove(saveFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal("error while removing prev temp file - " + err.Error())
	}

	err := SaveRequestsMade()
	if err != nil {
		t.Fatal("SaveRequestsMade returned error - " + err.Error())
	}

	if _, err := os.Stat(saveFilePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("SaveRequestsMade created a file when there are no records - " + err.Error())
	}
}

func TestLoadNoTempFile(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	if err := os.Remove(saveFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal("error while removing prev temp file - " + err.Error())
	}

	err := LoadRequestsMade()
	if err != nil {
		t.Fatal("LoadRequestsMade returned error - " + err.Error())
		return
	}

	if len(requestRecords) != 0 {
		t.Fatal("LoadRequestsMake managed to somehow extract records without a temp file")
	}
}