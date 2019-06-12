package util

import (
	"sync"
	"sync/atomic"
)

type Task func() error

func newWorker(p *Pool) *worker {
	w := &worker{
		tasks: make(chan Task),
		pool:  p,
	}
	w.run()
	return w
}

type worker struct {
	tasks chan Task
	pool  *Pool
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

func NewPool(cap int32) *Pool {
	p := &Pool{
		capacity: cap,
	}
	p.mux = &sync.Mutex{}
	p.cond = sync.NewCond(p.mux)
	p.once = &sync.Once{}
	return p
}

type Pool struct {
	capacity int32     // 最大运行数量
	running  int32     // 已运行数量
	workers  []*worker // 使用栈操作 slice
	mux      *sync.Mutex
	cond     *sync.Cond
	once     *sync.Once
}

// Put 把 task 放入 Pool 中, 应该是阻塞的?
func (p *Pool) Submit(task Task) {
	w := p.popWorker()
	w.tasks <- task
}

func (p *Pool) popWorker() *worker {
	p.mux.Lock()

	// 检测空闲 worker
	idx := len(p.workers) - 1
	if idx >= 0 {
		w := p.workers[idx]
		p.workers[idx] = nil
		p.workers = p.workers[:idx]
		p.running++

		p.mux.Unlock()
		return w
	}

	// 没有空闲的 worker
	// 检查已运行的数量
	if p.running < p.capacity { // 新建 worker
		w := newWorker(p)
		p.running++

		p.mux.Unlock()
		return w
	}
	p.mux.Unlock()

	// worker 数量达到上限
	// 条件变量
	p.cond.L.Lock()
	for !(p.running < p.capacity) {
		p.cond.Wait()
	}

	idx = len(p.workers) - 1
	w := p.workers[idx]
	p.workers[idx] = nil
	p.workers = p.workers[:idx]
	p.running++

	p.cond.L.Unlock()
	return w
}

func (p *Pool) pushWorker(w *worker) {
	if w == nil {
		return
	}
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	p.workers = append(p.workers, w)
	p.running--
	p.cond.Signal()
}

func (p *Pool) Running() int32 {
	return atomic.LoadInt32(&p.running)
}

func (p *Pool) Release() {
	// todo: 等待所有 worker 结束才退出
}
