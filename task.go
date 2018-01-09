package gmp

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	errNilFn   = errors.New("error: function is nil")
	errNotRun  = errors.New("error: pool is not running")
	errTimeout = errors.New("error: timed out")
)

// Task - task
type Task struct {
	ID       int64
	WorkerID int64
	Fn       func(...interface{}) interface{}
	Result   interface{}
	Args     []interface{}
	Error    error
}

// Add - add new task to pool
func (p *Pool) Add(fn func(...interface{}) interface{}, args ...interface{}) error {
	if fn == nil {
		return errNilFn
	}
	if !p.runningPool {
		return errNotRun
	}
	task := &Task{
		Fn:   fn,
		Args: args,
	}
	p.incNumTasks()
	task.ID = p.getNumTasks()
	p.queue.put(task)
	return nil
}

// SetTaskTimeout - set task timeout in second before send quit signal
func (p *Pool) SetTaskTimeout(t int) {
	p.quitTimeout = time.Duration(t) * time.Second
	p.useTimeout = true
}

func (p *Pool) exec(task *Task) *Task {
	defer func() {
		err := recover()
		if err != nil {
			task.Result = nil
			task.Error = fmt.Errorf("Recovery %v", err.(string))
		}
	}()
	if p.useTimeout {
		ch := make(chan interface{}, 1)
		// defer close(ch)

		go func() {
			ch <- task.Fn(task.Args...)
		}()

		select {
		case result := <-ch:
			task.Result = result
		case <-time.After(1 * time.Second):
			task.Error = errTimeout
		}
	} else {
		task.Result = task.Fn(task.Args...)
	}
	return task
}

func (p *Pool) getNumTasks() int64 {
	return atomic.LoadInt64(&p.numTasks)
}

func (p *Pool) incNumTasks() {
	atomic.AddInt64(&p.numTasks, 1)
}
