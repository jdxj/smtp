package util

import (
	"fmt"
	"testing"
)

func TestIdGen_GetID(t *testing.T) {
	for i := 0; i < 100000; i++ {
		str := IDGen.GetID()
		fmt.Println(str)
	}
}

func TestUint64ToBase62(t *testing.T) {
	str := Uint64ToBase62(2812840084295207690)
	fmt.Println(str)
}
