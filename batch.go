package forge

import "time"

// Batch will read values from the source channel and batch them in slices up
// to the specified size and sends them on the sink channel.
//
// It will finish batches within the specified timeout.
//
// If the source channel is closed the function will send the remaining batch
// and return.
//
// If the cancel channel is closed the function will return immediately. Data
// may be lost in this scenario.
func Batch(source <-chan Value, sink chan<- Value, cancel <-chan Signal, size int, timeout time.Duration) {
	// prepare slice
	var slice []Value

	// prepare timer
	var timer *time.Timer
	var trigger <-chan time.Time

	for {
		select {
		case value, ok := <-source:
			// check if source has been closed
			if !ok {
				// check remaining slice
				if len(slice) > 0 {
					// send slice
					select {
					case sink <- slice:
					case <-cancel:
					}
				}

				// stop timer if available
				if timer != nil {
					timer.Stop()
				}

				return
			}

			// create slice if nil
			if slice == nil {
				slice = make([]Value, 0, size)
			}

			// add value
			slice = append(slice, value)

			// set timer if missing
			if timer == nil && timeout > 0 {
				timer = time.NewTimer(timeout)
				trigger = timer.C
			}

			// check if slice is full
			if len(slice) == size {
				// send slice
				select {
				case sink <- slice:
				case <-cancel:
				}

				// make new slice
				slice = make([]Value, 0, size)

				// reset timer if available
				if timer != nil {
					timer.Stop()
					timer = nil
					trigger = nil
				}
			}
		case <-trigger:
			// send slice
			select {
			case sink <- slice:
			case <-cancel:
			}

			// reset slice
			slice = nil

			// reset timer
			timer.Stop()
			timer = nil
			trigger = nil

			continue
		case <-cancel:
			// stop timer if available
			if timer != nil {
				timer.Stop()
			}

			return
		}
	}
}