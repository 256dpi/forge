package forge

import (
	"errors"
	"sync"
)

// ErrStopped indicates that a supervisor has been stopped.
var ErrStopped = errors.New("stopped")

// ErrKilled indicates that a supervisor has been killed.
var ErrKilled = errors.New("killed")

// A Supervisor manages provides a stopping and killing mechanism.
type Supervisor struct {
	stopping chan Signal
	killed   chan Signal
	onceInit sync.Once
	onceStop sync.Once
	onceKill sync.Once
}

func (s *Supervisor) init() {
	// create the channels once
	s.onceInit.Do(func() {
		s.stopping = make(chan Signal)
		s.killed = make(chan Signal)
	})
}

// Stop will close the Stopping channel.
func (s *Supervisor) Stop() {
	s.init()

	// close channel once
	s.onceStop.Do(func() {
		close(s.stopping)
	})
}

// Stopping returns the channel closed by Stop.
func (s *Supervisor) Stopping() <-chan Signal {
	s.init()

	return s.stopping
}

// IsStopping returns whether Stop has been called.
func (s *Supervisor) IsStopping() bool {
	select {
	case <-s.stopping:
		return true
	default:
		return false
	}
}

// Kill will close the Stopping and Killed channel.
func (s *Supervisor) Kill() {
	s.init()

	// close channel once
	s.onceStop.Do(func() {
		close(s.stopping)
	})

	// close channel once
	s.onceKill.Do(func() {
		close(s.killed)
	})
}

// Killed returns the channel closed by Kill.
func (s *Supervisor) Killed() <-chan Signal {
	s.init()

	return s.killed
}

// IsKilled returns whether Kill has been called.
func (s *Supervisor) IsKilled() bool {
	select {
	case <-s.killed:
		return true
	default:
		return false
	}
}

// Status returns and error if Stop or Kill have been called.
func (s *Supervisor) Status() error {
	if s.IsKilled() {
		return ErrKilled
	} else if s.IsStopping() {
		return ErrStopped
	}

	return nil
}
