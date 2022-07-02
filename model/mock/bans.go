package mock

type Bans struct{}

func (b *Bans) Ban(int64) error {
	return nil
}

func (b *Bans) Unban(int64) error {
	return nil
}

func (b *Bans) List() ([]int64, error) {
	return []int64{}, nil
}

func (b *Bans) Banned(int64) (bool, error) {
	return false, nil
}
