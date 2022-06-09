package main

import (
	"fmt"
	"strings"
)

// whitelist is a list of group IDs where the bot works.
type whitelist struct {
	groupIDs map[int64]struct{}
}

// newWhitelist returns a new whitelist with the given group IDs.
func newWhitelist(groupIDs ...int64) *whitelist {
	w := &whitelist{make(map[int64]struct{})}
	for _, id := range groupIDs {
		w.groupIDs[id] = struct{}{}
	}
	return w
}

// add adds the given group ID to the whitelist.
func (w *whitelist) add(groupID int64) {
	w.groupIDs[groupID] = struct{}{}
}

// allow checks if the given group ID is in the whitelist.
func (w *whitelist) allow(groupID int64) bool {
	_, ok := w.groupIDs[groupID]
	return ok
}

func (w *whitelist) String() string {
	result := ""
	for id := range w.groupIDs {
		result = fmt.Sprintf("%s%d, ", result, id)
	}
	result = strings.TrimSuffix(result, ", ")
	return result
}
