package mock

type Users struct{}

func (u *Users) Insert(int64, int64) error {
	return nil
}
func (u *Users) Delete(int64, int64) error {
	return nil
}
func (u *Users) List(int64) ([]int64, error) {
	return []int64{}, nil
}
func (u *Users) Exists(int64, int64) (bool, error) {
	return false, nil
}
func (u *Users) Random(int64) (int64, error) {
	return 0, nil
}
func (u *Users) NRandom(int64, int) ([]int64, error) {
	return []int64{}, nil
}
