package module

import (
	"container/list"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	list *list.List
}

func (s *Store) Get() interface{} {
	var v interface{}
	s.mu.Lock()
	if s.list == nil {
		s.list = list.New()
	}
	e := s.list.Front()
	v = e.Value
	s.list.Remove(e)
	s.mu.Unlock()
	return v
}

func (s *Store) Add(v interface{}) {
	s.mu.Lock()
	if s.list == nil {
		s.list = list.New()
	}
	s.list.PushBack(v)
	s.mu.Unlock()
}
