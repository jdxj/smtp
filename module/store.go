package module

import (
	"container/list"
	"sync"
)

var Store *store

func init() {
	Store = &store{}
	Store.list = list.New()
}

type store struct {
	mu   sync.Mutex
	list *list.List
}

func (s *store) Get() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := s.list.Front()
	if e == nil {
		return nil
	}
	v := e.Value
	s.list.Remove(e)
	return v
}

func (s *store) Add(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list.PushBack(v)
}

func (s *store) Len() int {
	return s.list.Len()
}
