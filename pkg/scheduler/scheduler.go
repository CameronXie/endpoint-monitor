package scheduler

import (
	"context"
	"sync"
	"time"
)

type Task func(t time.Time)

type Scheduler interface {
	// Schedule runs a Task on any given interval.
	// Schedule also takes a context.Context, will stop running tasks on context cancellation.
	Schedule(ctx context.Context, interval time.Duration, f Task)

	// Wait blocks until all the tasks have been completed or gracefully closed.
	Wait()
}

type scheduler struct {
	wg sync.WaitGroup
}

func (s *scheduler) Schedule(ctx context.Context, interval time.Duration, task Task) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		t := time.NewTicker(interval)
		var wg sync.WaitGroup
		defer func() {
			t.Stop()
			wg.Wait()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case nt := <-t.C:
				wg.Add(1)
				go func(i time.Time) {
					defer wg.Done()
					task(i)
				}(nt)
			}
		}
	}()
}

func (s *scheduler) Wait() {
	s.wg.Wait()
}

func New() Scheduler {
	return new(scheduler)
}
