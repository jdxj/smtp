package util

import "sync"

type Task func() error

type worker struct {
	tasks chan Task
	pool *Pool
}

func (w *worker) run() {
	go func() {
		for task := range w.tasks {
			if task == nil {
				continue
			}

			task()
			// 运行完要放入池中
			w.pool.pushWorker(w)
		}
	}()
}

func NewPool(cap int) *Pool {
	p := &Pool{
		capacity:cap,
	}
	return p
}

type Pool struct {
	capacity int // 最大运行数量
	running int // 已运行数量
	workers []*worker
	mux sync.Mutex
}

// Put 把 task 放入 Pool 中, 应该是阻塞的?
func (p *Pool) Submit(task Task) {

}

func (p *Pool) popWorker() *worker {
	p.mux.Lock()
	defer p.mux.Unlock()

	idx := len(p.workers) -1
	if idx >= 0 { // 有正在等待的 worker

	}
	if len(p.workers) == 0 {
		return nil
	}
	//
	return nil
}

func (p *Pool) pushWorker(w *worker) {

}
