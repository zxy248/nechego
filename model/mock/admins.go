package mock

type Admins struct{}

func (a *Admins) Insert(int64) error {
	return nil
}
func (a *Admins) Delete(int64) error {
	return nil
}
func (a *Admins) List(int64) ([]int64, error) {
	return []int64{}, nil
}
func (a *Admins) Authorize(int64, int64) (bool, error) {
	return false, nil
}
func (a *Admins) InsertDaily(int64, int64) error {
	return nil
}
func (a *Admins) GetDaily(int64) (int64, error) {
	return 0, nil
}
func (a *Admins) DeleteDaily(int64) error {
	return nil
}
