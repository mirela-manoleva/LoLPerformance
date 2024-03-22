package limiter

import (
	"testing"
	"time"
)

func TestCanExecNoLimitsNoRecords(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	if !canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return true.\n" + getState())
	}
}

func TestCanExecNoRecords(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	AddLimit(5, time.Second)

	if !canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return true.\n" + getState())
	}
}

func TestCanExecNoLimits(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	addRecord(time.Now())

	if !canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return true.\n" + getState())
	}
}

func TestCanExec(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	limitResponse := 5

	AddLimit(limitResponse, time.Second)
	for i := 0; i < limitResponse-1; i++ {
		addRecord(time.Now())
	}

	if !canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return true.\n" + getState())
	}
}

func TestCannotExec(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	limitResponse := 5

	AddLimit(limitResponse, time.Second)
	for i := 0; i < limitResponse; i++ {
		addRecord(time.Now())
	}

	if canExecuteRequestNow() {
		t.Fatal("canExecuteRequestNow() should return false.\n" + getState())
	}
}

func TestRecordClear(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	limitResponse := 11
	limitPeriod := 50 * time.Millisecond

	AddLimit(limitResponse, limitPeriod)
	for i := 0; i < limitResponse; i++ {
		addRecord(time.Now())
	}

	time.Sleep(limitPeriod)
	canExecuteRequestNow()
	if len(requestRecords) != 0 {
		t.Fatal("canExecuteRequestNow() should have cleared the records.\n" + getState())
	}
}

// Check that canExecuteRequestNow() doesn't clear the records when it shouldn't
func TestRecordNotClear(t *testing.T) {
	clearAllRecords()
	ClearAllLimits()

	limitResponse := 5
	limitPeriod := 100 * time.Millisecond

	AddLimit(limitResponse, limitPeriod)
	for i := 0; i < limitResponse; i++ {
		addRecord(time.Now())
	}

	time.Sleep(limitPeriod - 50*time.Millisecond)
	numberOfRecords := len(requestRecords)

	canExecuteRequestNow()
	if len(requestRecords) != numberOfRecords {
		t.Fatal("canExecuteRequestNow() shouldn't have cleared any records." + getState())
	}
}
