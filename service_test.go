package forge

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceRunAndReporter(t *testing.T) {
	service := Service{}

	var err error
	service.Report(func(e error) {
		err = e
	})

	i := 0
	service.Run(1, func() error {
		i++
		if i == 2 {
			return ErrDone
		}
		return errors.New("foo")
	}, nil)

	<-service.Done()

	assert.Equal(t, 2, i)
	assert.Error(t, err)
}

func TestServicePipeline(t *testing.T) {
	service := &Service{}

	service.Report(func(error) {})

	i := 0
	service.Source(1, func(values chan<- Value) error {
		i++
		if i == 4 {
			return ErrDone
		}
		values <- i
		return nil
	}, 1)

	service.Batch(1, 1, time.Millisecond, 1)

	service.FilterFunc(1, func(v Value, out chan<- Value) error {
		out <- v
		return nil
	}, 1)

	var out []Value
	service.SinkFunc(1, func(v Value) error {
		out = append(out, v)
		return nil
	})

	<-service.Done()

	assert.Equal(t, []Value{[]Value{1}, []Value{2}, []Value{3}}, out)
}

func TestServiceSendAndStop(t *testing.T) {
	service := &Service{}

	service.Report(func(error) {})

	service.Open(0)

	var out1 []string
	service.Sink(1, func(in <-chan Value) error {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					out1 = append(out1, "done")
					return ErrDone
				}

				out1 = append(out1, "ok")
				time.Sleep(10 * time.Millisecond)
			}
		}
	})

	var out2 []string
	go func() {
		for i := 1; i <= 5; i++ {
			err := service.Send(i)
			if err != nil {
				out2 = append(out2, "error")
			} else {
				out2 = append(out2, "ok")
			}
		}
	}()

	time.Sleep(25 * time.Millisecond)
	service.Stop()
	<-service.Done()

	assert.Equal(t, []string{"ok", "ok", "ok", "done"}, out1)
	assert.Equal(t, []string{"ok", "ok", "ok", "error", "error"}, out2)
	assert.Equal(t, ErrStopped, service.Status())
}
