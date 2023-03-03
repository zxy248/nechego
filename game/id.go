package game

// IDer is implemented by any value that has a unique ID.
type IDer interface {
	ID() int64
}

// ID is a trivial implementation of the IDer interface.
type ID int64

func (id ID) ID() int64 {
	return int64(id)
}
