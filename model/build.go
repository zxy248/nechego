package model

import "strings"

func concat(elems ...string) string {
	return strings.Join(elems, " ")
}

func and(elems ...string) string {
	return strings.Join(append([]string{"1 = 1"}, elems...), " and ")
}
