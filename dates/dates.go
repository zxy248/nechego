package dates

import "time"

func Today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}
