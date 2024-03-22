package limiter

/*
	File description:
	Defines the limit struct and gives the user the ability to modify rules.
*/

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
*/
type limit struct {
	requestCount int
	period       time.Duration
}

/*
Helper function to reduce code size.
*/
func (main limit) isEqual(other limit) bool {
	return main.period == other.period && main.requestCount == other.requestCount
}

/*
Adds a limit. Calls to the SendRequest() will abide to the limit.
*/
func AddLimit(requestCount int, period time.Duration) {
	limit := limit{requestCount, period}
	i := 0
	for ; i < len(requestLimits) && limit.period > requestLimits[i].period; i++ {
	}

	// Do not store copies
	if i < len(requestLimits) && limit.isEqual(requestLimits[i]) {
		return
	}

	requestLimits = slices.Insert(requestLimits, i, limit)
}

/*
Removes previously declared limit.
*/
func RemoveLimit(lim limit) {
	for i := 0; i < len(requestLimits); i++ {
		if lim.isEqual(requestLimits[i]) {
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
