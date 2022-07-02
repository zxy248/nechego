package mock

type Eblans struct{}

func (e *Eblans) Insert(int64, int64) error {
	return nil
}

func (e *Eblans) Get(int64) (int64, error) {
	return 0, nil
}

func (e *Eblans) Delete(int64) error {
	return nil
}
