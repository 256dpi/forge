package forge

import "sync"

// Run will launch multiple goroutines that execute the specified task. If a
// finalizer is configured, it will be called once all tasks returned.
func Run(n int, task func(), finalizer func()) {
	// prepare wait group
	var wg sync.WaitGroup
	wg.Add(n)

	// run tasks
	for i := 0; i < n; i++ {
		go func() {
			task()
			wg.Done()
		}()
	}

	// run finalizer if available
	if finalizer != nil {
		go func() {
			wg.Wait()
			finalizer()
		}()
	}
}
