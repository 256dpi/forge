package forge

import (
	"sync"
	"time"
)

// Service wraps Pipeline with a Terminator and Reporter.
type Service struct {
	Pipeline
	Terminator
	Reporter

	mutex sync.Mutex
	once  sync.Once
}

// Run will run the specified task.
func (s *Service) Run(n int, task func() error, finalizer func()) {
	s.Manager.Run(n, func() {
		s.Repeat(task)
	}, finalizer)
}

// Send will send the specified value to the pipeline by respecting an eventual
// closing of the service.
func (s *Service) Send(v Value) error {
	// acquire mutex
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// check status
	err := s.Status()
	if err != nil {
		return err
	}

	// queue value
	select {
	case s.Input() <- v:
	case <-s.Stopping():
		return s.Status()
	}

	return nil
}

// Source will run the source task that fills the provided channel with values.
func (s *Service) Source(n int, fn func(chan<- Value) error, buffer int) {
	s.Pipeline.Source(n, func(out chan<- Value) {
		s.Repeat(func() error {
			return fn(out)
		})
	}, buffer)
}

// Filter is an intermediary task that processes values.
func (s *Service) Filter(n int, fn func(<-chan Value, chan<- Value) error, buffer int) {
	s.Pipeline.Filter(n, func(in <-chan Value, out chan<- Value) {
		s.Repeat(func() error {
			return fn(in, out)
		})
	}, buffer)
}

// FilterFunc augments Filter by running the specified function once per
// received value.
func (s *Service) FilterFunc(n int, fn func(Value, chan<- Value) error, buffer int) {
	s.Filter(n, func(in <-chan Value, out chan<- Value) error {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return ErrDone
				}

				// call function
				err := fn(v, out)
				if err != nil {
					return err
				}
			case <-s.Killed():
				return ErrDone
			}
		}
	}, buffer)
}

// Batch is an intermediary task that batches values up.
func (s *Service) Batch(n int, sizer func(Value) int, limit int, timeout time.Duration, buffer int) {
	s.Pipeline.Batch(n, sizer, limit, s.Killed(), timeout, buffer)
}

// Sink is the final task that receives all processed values.
func (s *Service) Sink(n int, fn func(<-chan Value) error) {
	s.Pipeline.Sink(n, func(in <-chan Value) {
		s.Repeat(func() error {
			return fn(in)
		})
	})
}

// SinkFunc augments Sink by running the specified function once per received
// value.
func (s *Service) SinkFunc(n int, fn func(Value) error) {
	s.Sink(n, func(in <-chan Value) error {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return ErrDone
				}

				// call function
				err := fn(v)
				if err != nil {
					return err
				}
			case <-s.Killed():
				return ErrDone
			}
		}
	})
}

// Stop will stop the Terminator and close the input channel if the Pipeline has
// been opened.
func (s *Service) Stop() {
	// forward call (will unblock blocked Send calls)
	s.Terminator.Stop()

	// acquire mutex
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// close input once
	s.once.Do(func() {
		if s.Input() != nil {
			close(s.Input())
		}
	})
}
