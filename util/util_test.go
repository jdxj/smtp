package util

import (
	"fmt"
	"testing"
	"time"
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

func TestGoroutinePool(t *testing.T) {
	p := NewPool(1000)

	task1 := func() error {
		fmt.Println("I'm task1! working...")
		// 模拟工作中
		time.Sleep(5 * time.Second)
		fmt.Println("I'm task1! work end!")
		return nil
	}

	task2 := func() error {
		fmt.Println("I'm task2! working...")
		// 模拟工作中
		time.Sleep(10 * time.Second)
		fmt.Println("I'm task2! work end!")
		return nil
	}

	task3 := func() error {
		fmt.Println("I'm task3! working...")
		// 模拟工作中
		time.Sleep(15 * time.Second)
		fmt.Println("I'm task3! work end!")
		return nil
	}

	tasks := []Task{task1, task2, task3}
	// 自动填装任务
	go func() {
		for i := 0; i != 10; i++ {
			task := tasks[i%3]
			fmt.Println(i, " Submitting...")
			p.Submit(task)
			fmt.Println(i, " Submit success!")

			// 填装间隔
			time.Sleep(500 * time.Microsecond)
			fmt.Println(i, "pool running: ", p.Running())
			fmt.Println("--------------------------")
		}
	}()

	// 注协程等待
	time.Sleep(20 * time.Minute)
}

func BenchmarkGoroutinePool(b *testing.B) {
	f := func() error {
		time.Sleep(time.Second)
		return nil
	}
	p := NewPool(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.Submit(f)
	}
}

func TestChan(t *testing.T) {
	c := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			c <- i
			time.Sleep(time.Second)
		}
		close(c)
	}()

	for {
		v, ok := <-c
		if ok {
			fmt.Println(v)
		} else {
			fmt.Println("end!")
			break
		}
	}
}

func TestNilChan(t *testing.T) {
	var c chan int
	go func() {
		for v := range c {
			fmt.Println("read v: ", v)
		}
	}()

	close(c)
	time.Sleep(10 * time.Second)
	fmt.Println("main goroutine!")
}
