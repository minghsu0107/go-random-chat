package common

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

var numCPU int

func init() {
	numCPU = runtime.NumCPU()
}

// TaskResult represent result of task
type TaskResult struct {
	Result interface{}
	Err    error
}

// Task represent a task
type Task struct {
	ctx      context.Context
	executor func(context.Context) (interface{}, error)
	future   chan *TaskResult
}

// NewTask create new task
func NewTask(ctx context.Context, executor func(context.Context) (interface{}, error)) *Task {
	return &Task{
		ctx:      ctx,
		executor: executor,
		future:   make(chan *TaskResult, 1),
	}
}

// Execute task
func (t *Task) Execute() {
	var result interface{}
	var err error

	if t.executor != nil {
		result, err = t.executor(t.ctx)
	}

	t.future <- &TaskResult{Result: result, Err: err}
}

// Result pushed via channel
func (t *Task) Result() <-chan *TaskResult {
	return t.future
}

// Option represents pool option
type Option struct {
	// NumberWorker number of workers
	// Default: runtime.NumCPU()
	NumberWorker int
}

func (o *Option) normalize() {
	if o.NumberWorker <= 0 {
		o.NumberWorker = numCPU
	}
}

// Pool is a lightweight worker pool with capable of auto-expand on demand
type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc

	opt Option

	wg        sync.WaitGroup
	taskQueue chan *Task

	state uint32 // 0: not start, 1: started, 2: stopped
}

// NewPool create new worker pool
func NewPool(ctx context.Context, opt Option) (p *Pool) {
	if ctx == nil {
		ctx = context.Background()
	}

	// normalize option
	opt.normalize()

	// set up pool
	p = &Pool{
		opt:       opt,
		taskQueue: make(chan *Task, opt.NumberWorker),
	}
	p.ctx, p.cancel = context.WithCancel(ctx)
	return
}

// Start workers
func (p *Pool) Start() {
	if atomic.CompareAndSwapUint32(&p.state, 0, 1) {
		numWorker := p.opt.NumberWorker

		p.wg.Add(numWorker)
		for i := 0; i < numWorker; i++ {
			go p.startWorker()
		}
	}
}

func (p *Pool) startWorker() {
	for task := range p.taskQueue {
		task.Execute()
	}
	p.wg.Done()
}

// Stop worker and wait all task done
func (p *Pool) Stop() {
	if atomic.CompareAndSwapUint32(&p.state, 1, 2) || atomic.CompareAndSwapUint32(&p.state, 0, 2) {
		// cancel context
		p.cancel()

		// wait child workers
		close(p.taskQueue)
		p.wg.Wait()
	}
}

// Do a task
func (p *Pool) Do(t *Task) {
	if t != nil {
		p.push(t)
	}
}

func (p *Pool) push(t *Task) {
	select {
	case <-p.ctx.Done():
	case p.taskQueue <- t:
	}
}
