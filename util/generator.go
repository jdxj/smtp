package util

import (
	"math/rand"
	"sync"
	"time"
)

const printable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var IDGen *idGen

func init() {
	IDGen = &idGen{}
}

type idGen struct {
	mu  sync.Mutex
	inc uint16
	//tim uint16
	//ran uint32
	// inc-ran-tim
}

// GetID 返回自增 id.
func (ig *idGen) GetID() string {
	var id uint64
	tim := uint16(time.Now().Unix())
	ran := rand.Uint32()

	ig.mu.Lock()
	inc := ig.inc

	id = id | uint64(inc)
	id = id << 32

	id = id | uint64(ran)
	id = id << 16

	id = id | uint64(tim)

	ig.inc++
	ig.mu.Unlock()

	return Uint64ToBase62(id)
}

func Uint64ToBase62(num uint64) (res string) {
	for i := 0; i < 8; i++ {
		idx := byte(num) % 62
		res += string(printable[idx])

		num = num >> 8
	}
	return res
}
