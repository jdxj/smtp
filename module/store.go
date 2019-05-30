package module

import (
	"sync"
	"time"
)

var Store *store

func init() {
	Store = &store{
		M: &sync.Map{},
	}
}

type store struct {
	M *sync.Map
}

func (s *store) DelMail(dur time.Duration, key string) {
	f := func() {
		s.M.Delete(key)
	}
	time.AfterFunc(dur, f)
}

func (s *store) DelUser(dur time.Duration, key string) {
	s.DelMail(dur, key)
}
