package forge

import (
	"sync"
)

// A Manager manages multiple finite running tasks.
type Manager struct {
	done  chan struct{}
	mutex sync.Mutex
}

// Run will run the specified task.
func (m *Manager) Run(n int, task func(), finalizer func()) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// get previous done channel
	previousDone := m.done

	// create new done channel
	currentDone := make(chan struct{})

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
func (m *Manager) Done() <-chan struct{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.done
}

// Finished will return whether all ran tasks have returned.
func (m *Manager) Finished() bool {
	select {
	case <-m.Done():
		return true
	default:
		return false
	}
}
