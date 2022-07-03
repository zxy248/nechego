package sqlite

import "time"

type Messages struct {
	DB *DB
}

// TODO: UserCount returns the number of messages sent by the user since the specified time.
func (m *Messages) UserCount(gid, uid int64, since time.Time) (int, error) {
	return 0, nil
}

// TODO: TotalCount returns the number of messages sent total since the specified time.
func (m *Messages) TotalCount(gid int64, since time.Time) (int, error) {
	return 0, nil
}
