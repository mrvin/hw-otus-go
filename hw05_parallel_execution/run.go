package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	// ErrErrorsLimitExceeded is errors limit exceeded.
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	// ErrIncorrectParam is errors incorrect parameters.
	ErrIncorrectParam = errors.New("errors incorrect parameters")
)

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, N int, M int) error { //nolint:gocritic,cyclop
	taskCh := make(chan Task)
	var errCnt int32

	if N < 1 || M < 0 || len(tasks) == 0 {
		return ErrIncorrectParam
	}

	numConsumer := N
	if len(tasks) < N {
		numConsumer = len(tasks)
	}

	wg := sync.WaitGroup{}
	wg.Add(numConsumer)
	for i := 0; i < numConsumer; i++ {
		go func() {
			defer wg.Done()

			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errCnt, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCnt) >= int32(M) {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if errCnt >= int32(M) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
