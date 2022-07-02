package mock

import "nechego/input"

type Forbid struct{}

func (f *Forbid) Forbid(int64, input.Command) error {
	return nil
}
func (f *Forbid) Permit(int64, input.Command) error {
	return nil
}
func (f *Forbid) Forbidden(int64, input.Command) (bool, error) {
	return false, nil
}
func (f *Forbid) List(int64) ([]input.Command, error) {
	return []input.Command{}, nil
}
