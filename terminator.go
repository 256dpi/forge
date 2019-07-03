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
	stopping chan struct{}
	killed   chan struct{}

	onceInit sync.Once
	onceStop sync.Once
	onceKill sync.Once

	notifiers []func()
}

func (t *Terminator) init() {
	// create the channels once
	t.onceInit.Do(func() {
		t.stopping = make(chan struct{})
		t.killed = make(chan struct{})
	})
}

// Notify will store the specified callback and call it once the terminator has
// been stopped.
func (t *Terminator) Notify(fn func()) {
	t.init()

	// add tracker
	t.notifiers = append(t.notifiers, fn)
}

// Stop will close the Stopping channel.
func (t *Terminator) Stop() {
	t.init()

	// close channel once
	t.onceStop.Do(func() {
		close(t.stopping)

		// call notifiers
		go func() {
			for _, t := range t.notifiers {
				t()
			}
		}()
	})
}

// Stopping returns the channel closed by Stop.
func (t *Terminator) Stopping() <-chan struct{} {
	t.init()

	return t.stopping
}

// IsStopping returns whether Stop has been called.
func (t *Terminator) IsStopping() bool {
	select {
	case <-t.stopping:
		return true
	default:
		return false
	}
}

// Kill will close the Stopping and Killed channel.
func (t *Terminator) Kill() {
	t.init()

	// close channel once
	t.onceStop.Do(func() {
		close(t.stopping)
	})

	// close channel once
	t.onceKill.Do(func() {
		close(t.killed)
	})
}

// Killed returns the channel closed by Kill.
func (t *Terminator) Killed() <-chan struct{} {
	t.init()

	return t.killed
}

// IsKilled returns whether Kill has been called.
func (t *Terminator) IsKilled() bool {
	select {
	case <-t.killed:
		return true
	default:
		return false
	}
}

// Status returns and error if Stop or Kill have been called.
func (t *Terminator) Status() error {
	if t.IsKilled() {
		return ErrKilled
	} else if t.IsStopping() {
		return ErrStopped
	}

	return nil
}
