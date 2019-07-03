package forge

import "sync/atomic"

// Run will launch multiple goroutines that execute the specified task. If a
// finalizer is configured, it will be called once all tasks returned.
func Run(n int, task func(), finalizer func()) {
	// prepare wait group
	counter := int64(n)

	// run tasks
	for i := 0; i < n; i++ {
		go func() {
			task()

			// run finalizer if available if task is last to return
			if atomic.AddInt64(&counter, -1) == 0 {
				if finalizer != nil {
					finalizer()
				}
			}
		}()
	}
}
