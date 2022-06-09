package main

import (
	"fmt"
	"strings"
)

// owners is a list of user IDs that have access to the administrative functions.
type owners struct {
	userIDs map[int64]struct{}
}

// newOwners returns a new owners with the given user IDs.
func newOwners(userIDs ...int64) *owners {
	o := &owners{make(map[int64]struct{})}
	for _, id := range userIDs {
		o.userIDs[id] = struct{}{}
	}
	return o
}

// add adds the given user ID to the list of owners.
func (o *owners) add(userID int64) {
	o.userIDs[userID] = struct{}{}
}

// check checks if the given user ID is in the list of owners.
func (o *owners) check(userID int64) bool {
	_, ok := o.userIDs[userID]
	return ok
}

func (o *owners) String() string {
	result := ""
	for id := range o.userIDs {
		result = fmt.Sprintf("%s%d, ", result, id)
	}
	result = strings.TrimSuffix(result, ", ")
	return result
}
