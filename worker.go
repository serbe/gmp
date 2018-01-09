package gmp

func (p *Pool) worker(id int64) {
	for {
		task, ok := p.queue.get()
		if ok {
			task.WorkerID = id
			p.result.put(p.exec(task))
		} else {
			break
		}
	}
}

// RunWorkers - run workers
// use it after add all tasks
func (p *Pool) RunWorkers() {
	var i int64
	for i = 0; i < p.numWorkers; i++ {
		go p.worker(i)
	}
	p.runningPool = false
}
