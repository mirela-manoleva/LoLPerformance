package limiter

import (
	"slices"
	"time"
)

/*
	Keeps track of all user-defined limits. Can be edited with AddLimit, RemoveLimit and ClearAllLimits.
	It is intended to be a set of sorted limits.
	Manually appending or modifying it outside of the 3 intended functions will result in undefined behaviour.
*/
var requestLimits []limit

/*
	Type to encapsulate the definition of an API Limit.
	A limit is defined as requestsAllowed per period of time.
 */
type limit struct {
	requestsAllowed		int
	period						time.Duration
}

/*
	Helper function to reduce code size.
*/
func (main limit) isEqual(other limit) bool {
	return main.period == other.period && main.requestsAllowed == other.requestsAllowed
}

/*
	Adds a limit. Calls to the SendRegulatedRequest() will abide to the limit.
*/
func AddLimit(allowedRequests int, timeRange time.Duration) {
	limit := limit{requestsAllowed: allowedRequests, period: timeRange}
	i := 0
	for ; i < len(requestLimits) && limit.period > requestLimits[i].period; i++ {}

	// Do not store copies
	if i < len(requestLimits) && limit.isEqual(requestLimits[i]) {
		return
	}

	requestLimits = slices.Insert(requestLimits, i, limit)
}

/*
	Removes previously declared limit.
*/
func RemoveLimit(limit limit) {
	for i := 0; i < len(requestLimits); i++ {
		if limit.isEqual(requestLimits[i]) {
			requestLimits = slices.Delete(requestLimits, i, i+1)
			return
		}
	}
}

/*
	Removes all previously declared limits.
*/
func ClearAllLimits() {
	requestLimits = nil
}

/*
	Returns false if the request will break any limit in requestLimits.
*/
func canExecuteRequestNow() bool {
	if len(requestLimits) == 0 {
		return true;
	}

	currentTime := time.Now()
	counter := 0

	for i, j := len(requestRecords) - 1, 0; i >= 0; i-- {
		counter++

		if counter >= requestLimits[j].requestsAllowed {
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