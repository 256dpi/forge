package forge

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceRunAndReporter(t *testing.T) {
	s := Service{}

	var err error
	s.Report(func(e error) {
		err = e
	})

	i := 0
	s.Run(1, func() error {
		i++
		if i == 2 {
			return ErrDone
		}
		return errors.New("foo")
	}, nil)

	<-s.Done()

	assert.Equal(t, 2, i)
	assert.Error(t, err)
}

func TestServicePipeline(t *testing.T) {
	s := &Service{}

	s.Report(func(error) {})

	i := 0
	s.Source(1, func(values chan<- Value) error {
		i++
		if i == 4 {
			return ErrDone
		}
		values <- i
		return nil
	}, 1)

	s.Batch(1, 1, time.Millisecond, 1)

	s.FilterFunc(1, func(v Value, out chan<- Value) error {
		out <- v
		return nil
	}, 1)

	var out []Value
	s.SinkFunc(1, func(v Value) error {
		out = append(out, v)
		return nil
	})

	<-s.Done()

	assert.Equal(t, []Value{[]Value{1}, []Value{2}, []Value{3}}, out)
}

func TestServiceSendAndStop(t *testing.T) {
	s := &Service{}

	s.Report(func(error) {})

	s.Open(0)

	var ret1 []string
	s.Sink(1, func(in <-chan Value) error {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					ret1 = append(ret1, "done")
					return ErrDone
				}

				ret1 = append(ret1, "ok")
				time.Sleep(10 * time.Millisecond)
			}
		}
	})

	var ret2 []string
	go func() {
		for i := 1; i <= 5; i++ {
			err := s.Send(i)
			if err != nil {
				ret2 = append(ret2, "error")
			} else {
				ret2 = append(ret2, "ok")
			}
		}
	}()

	time.Sleep(25 * time.Millisecond)
	s.Stop()
	<-s.Done()

	assert.Equal(t, []string{"ok", "ok", "ok", "done"}, ret1)
	assert.Equal(t, []string{"ok", "ok", "ok", "error", "error"}, ret2)
	assert.Equal(t, ErrStopped, s.Status())
}
