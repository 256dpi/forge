package forge

import (
	"sync"
)

// A Manager manages multiple running tasks.
type Manager struct {
	done  chan Signal
	mutex sync.Mutex
}

// Run will run the specified task.
func (m *Manager) Run(n int, task func(), finalizer func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// get previous done channel
	previousDone := m.done

	// create new done channel
	currentDone := make(chan Signal)

	// run task
	Run(n, task, func() {
		// call finalizer if available
		if finalizer != nil {
			finalizer()
		}

		// wait for previous done if available
		if previousDone != nil {
			<-previousDone
		}

		// close current done
		close(currentDone)
	})

	// set done
	m.done = currentDone
}

// Done will return a channel that is closed once all until now started tasks
// have returned.
func (m *Manager) Done() <-chan Signal {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.done
}

// Finished will return whether all ran tasks have returned.
func (m *Manager) Finished() bool {
	select {
	case <-m.done:
		return true
	default:
		return false
	}
}
