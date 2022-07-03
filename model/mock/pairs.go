package mock

type Pairs struct{}

func (p *Pairs) Insert(int64, int64, int64) error {
	return nil
}
func (p *Pairs) Get(int64) (int64, int64, error) {
	return 0, 0, nil
}
