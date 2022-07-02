package mock

type Whitelist struct{}

func (w *Whitelist) Insert(int64) error {
	return nil
}

func (w *Whitelist) Allow(int64) (bool, error) {
	return false, nil
}
