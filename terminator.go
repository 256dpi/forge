package forge

import (
	"errors"
	"sync"
)

// ErrStopped indicates that a Terminator has been stopped.
var ErrStopped = errors.New("stopped")

// ErrKilled indicates that a Terminator has been killed.
var ErrKilled = errors.New("killed")

// Terminator provides a stopping and killing mechanism.
type Terminator struct {
	stopping chan Signal
	killed   chan Signal
	onceInit sync.Once
	onceStop sync.Once
	onceKill sync.Once
}

func (s *Terminator) init() {
	// create the channels once
	s.onceInit.Do(func() {
		s.stopping = make(chan Signal)
		s.killed = make(chan Signal)
	})
}

// Stop will close the Stopping channel.
func (s *Terminator) Stop() {
	s.init()

	// close channel once
	s.onceStop.Do(func() {
		close(s.stopping)
	})
}

// Stopping returns the channel closed by Stop.
func (s *Terminator) Stopping() <-chan Signal {
	s.init()

	return s.stopping
}

// IsStopping returns whether Stop has been called.
func (s *Terminator) IsStopping() bool {
	select {
	case <-s.stopping:
		return true
	default:
		return false
	}
}

// Kill will close the Stopping and Killed channel.
func (s *Terminator) Kill() {
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
func (s *Terminator) Killed() <-chan Signal {
	s.init()

	return s.killed
}

// IsKilled returns whether Kill has been called.
func (s *Terminator) IsKilled() bool {
	select {
	case <-s.killed:
		return true
	default:
		return false
	}
}

// Status returns and error if Stop or Kill have been called.
func (s *Terminator) Status() error {
	if s.IsKilled() {
		return ErrKilled
	} else if s.IsStopping() {
		return ErrStopped
	}

	return nil
}
