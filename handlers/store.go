package handlers

import "sync"

type Store[T any] struct {
	data sync.Map
}

func (s *Store[T]) Get(key int64, new T) (value T, done func()) {
	v, _ := s.data.LoadOrStore(key, &exclusive[T]{value: new})
	x := (v).(*exclusive[T])
	x.mu.Lock()
	return x.value, x.mu.Unlock
}

type exclusive[T any] struct {
	value T
	mu    sync.Mutex
}
