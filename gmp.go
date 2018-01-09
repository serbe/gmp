package gmp

import (
	"time"
)

// Pool - specification of gopool
type Pool struct {
	useTimeout  bool
	runningPool bool
	numWorkers  int64
	freeWorkers int64
	numTasks    int64
	queue       *taskList
	result      *taskList
	quitTimeout time.Duration
}

// New - create new gorourine pool
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.queue = new(taskList)
	p.result = new(taskList)
	p.runningPool = true
	return p
}
