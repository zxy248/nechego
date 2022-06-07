package main

import (
	"fmt"
	"strings"
)

type whitelist struct {
	groupIDs map[int64]struct{}
}

func newWhitelist(groupIDs ...int64) *whitelist {
	w := &whitelist{make(map[int64]struct{})}
	for _, id := range groupIDs {
		w.groupIDs[id] = struct{}{}
	}
	return w
}

func (w *whitelist) add(groupID int64) {
	w.groupIDs[groupID] = struct{}{}
}

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
