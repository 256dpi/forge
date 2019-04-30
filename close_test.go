package forge

import (
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	s1 := &Service{}
	s1.Report(func(error) {})
	s1.Run(1, func() error {
		return ErrDone
	}, nil)

	s2 := &Service{}
	s2.Report(func(error) {})
	s2.Run(1, func() error {
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, s1, s2)
}

func TestCloseStopTimeout(t *testing.T) {
	s1 := &Service{}
	s1.Report(func(error) {})
	s1.Run(1, func() error {
		<-s1.Killed()
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, s1)
}

func TestCloseKillTimeout(t *testing.T) {
	done := make(chan Signal)

	s1 := &Service{}
	s1.Report(func(error) {})
	s1.Run(1, func() error {
		<-done
		return ErrDone
	}, nil)

	Close(time.Millisecond, time.Millisecond, s1)

	close(done)
}
