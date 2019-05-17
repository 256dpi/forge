package forge

import (
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	service1 := &Service{}
	service1.Report(func(error) {})
	service1.Run(1, func() error {
		return ErrDone
	}, nil)

	service2 := &Service{}
	service2.Report(func(error) {})
	service2.Run(1, func() error {
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, service1, service2)
}

func TestCloseStopTimeout(t *testing.T) {
	service := &Service{}
	service.Report(func(error) {})
	service.Run(1, func() error {
		<-service.Killed()
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, service)
}

func TestCloseKillTimeout(t *testing.T) {
	done := make(chan Signal)

	service := &Service{}
	service.Report(func(error) {})
	service.Run(1, func() error {
		<-done
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, service)

	close(done)
}
