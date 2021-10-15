package scheduler

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_Schedule(t *testing.T) {
	a := assert.New(t)

	s := New()
	ctx, cancel := context.WithCancel(context.Background())
	interval, numTest, variance := time.Duration(10), 10, float64(2)
	ts := &concurrentTimeSlice{slice: make([]time.Time, 0)}

	s.Schedule(ctx, interval*time.Millisecond, func(i time.Time) {
		time.Sleep(interval * time.Millisecond)
		ts.Append(i)
	})

	time.AfterFunc(time.Duration(numTest)*interval*time.Millisecond, func() {
		cancel()
	})

	s.Wait()

	var head time.Time
	a.Equal(numTest, len(ts.toSlice()))
	for _, v := range ts.toSlice() {
		if head.Nanosecond() == 0 {
			head = v
			continue
		}

		a.LessOrEqual(math.Abs(float64(interval)-float64(v.Sub(head).Milliseconds())), variance)
		head = v
	}
}

type concurrentTimeSlice struct {
	sync.RWMutex
	slice []time.Time
}

func (s *concurrentTimeSlice) Append(t time.Time) {
	s.Lock()
	defer s.Unlock()
	s.slice = append(s.slice, t)
}

func (s *concurrentTimeSlice) toSlice() []time.Time {
	s.RLock()
	defer s.RUnlock()
	return s.slice
}
