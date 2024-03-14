package data

import "time"

type ChatData struct {
	Admin     int64     `json:"admin"`
	Eblan     int64     `json:"eblan"`
	Pair1     int64     `json:"pair1"`
	Pair2     int64     `json:"pair2"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
}
