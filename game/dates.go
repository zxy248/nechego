package game

import "time"

func today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func tomorrow() time.Time {
	return today().AddDate(0, 0, 1)
}
