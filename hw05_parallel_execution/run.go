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
	errLimit int
	wg       sync.WaitGroup
	mu       sync.Mutex
	tasks    []Task
	state    runnerState
	errCount int
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
	r.mu.Lock()
	defer r.mu.Unlock()

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

	freeWorkers, busyWorkers := delegateTasks(n)

	doTasks(freeWorkers, busyWorkers)

	if r.state == failed {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func delegateTasks(workersCount int) (chan<- chan Task, <-chan chan Task) {
	busyWorkers := make(chan chan Task, workersCount)

	freeWorkers := make(chan chan Task, workersCount)
	for i := 0; i < workersCount; i++ {
		freeWorkers <- make(chan Task, 1)
	}

	go func() {
		for worker := range freeWorkers {
			r.syncState()

			if r.state == failed || r.state == idle {
				close(busyWorkers)
				return
			}

			worker <- r.findTask()
			busyWorkers <- worker
		}
	}()

	return freeWorkers, busyWorkers
}

func doTasks(freeWorkers chan<- chan Task, busyWorkers <-chan chan Task) {
	for worker := range busyWorkers {
		worker := worker
		t := <-worker
		r.wg.Add(1)

		go func() {
			err := t()
			if err != nil {
				r.addErr()
			}

			r.wg.Done()
			freeWorkers <- worker
		}()
	}

	r.wg.Wait()
}
