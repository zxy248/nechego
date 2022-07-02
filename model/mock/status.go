package mock

type Status struct{}

func (s *Status) Enable(int64) error {
	return nil
}
func (s *Status) Active(int64) (bool, error) {
	return false, nil
}
func (s *Status) Disable(int64) error {
	return nil
}
