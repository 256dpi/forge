package forge

import "time"

// Closer is a closable service or struct that embeds a service.
type Closer interface {
	Done() <-chan struct{}
	Stop()
	Kill()
}

// Close will stop and kill the provided closers in the specified timeouts.
func Close(stop, kill time.Duration, closers ...Closer) bool {
	// create timeout channel
	tch := time.After(stop)

	// stop all closers
	for _, c := range closers {
		c.Stop()
	}

	// await all closers, or timeout
	for _, c := range closers {
		select {
		case <-c.Done():
			continue
		case <-tch:
			break
		}

		break
	}

	// reset timeout channel
	tch = time.After(kill)

	// kill all closers
	for _, c := range closers {
		c.Kill()
	}

	// await all closers, or timeout
	for _, c := range closers {
		select {
		case <-c.Done():
			// continue
		case <-tch:
			return false
		}
	}

	return true
}
