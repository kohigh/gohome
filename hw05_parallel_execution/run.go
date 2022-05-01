package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type runnerState int

const (
	processing runnerState = iota
	failed
	idle
)

type runner struct {
	errLimit    int
	tasks       []Task
	busyWorkers chan chan Task
	freeWorkers chan chan Task
	wg          sync.WaitGroup
	mu          sync.Mutex
	state       runnerState
	errCount    int
}

func (r *runner) syncState() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.errCount > r.errLimit {
		r.state = failed
	}

	if len(r.tasks) == 0 {
		r.state = idle
	}
}

func (r *runner) addErr() {
	r.mu.Lock()
	r.errCount++
	r.mu.Unlock()
}

func (r *runner) findTask() Task {
	task := r.tasks[len(r.tasks)-1]
	r.tasks = r.tasks[:len(r.tasks)-1]

	return task
}

var r runner

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 || len(tasks) == 0 || m < 1 {
		return nil
	}

	r = runner{state: processing, errLimit: m, tasks: tasks}

	allocWorkers(n)

	go delegateTasks()

	doTasks()

	if r.state == failed {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func allocWorkers(workersCount int) {
	r.busyWorkers = make(chan chan Task, workersCount)
	r.freeWorkers = make(chan chan Task, workersCount)

	for i := 0; i < workersCount; i++ {
		r.freeWorkers <- make(chan Task, 1)
	}
}

func delegateTasks() {
	for worker := range r.freeWorkers {
		r.syncState()

		if r.state == failed || r.state == idle {
			close(r.busyWorkers)
			return
		}

		worker <- r.findTask()
		r.busyWorkers <- worker
	}
}

func doTasks() {
	for worker := range r.busyWorkers {
		worker := worker
		t := <-worker
		r.wg.Add(1)

		go func() {
			err := t()
			if err != nil {
				r.addErr()
			}

			r.wg.Done()
			r.freeWorkers <- worker
		}()
	}

	r.wg.Wait()
}
