package forge

import "time"

// Batch will read values from the source channel and batch them in slices up
// to the specified limit and sends them on the sink channel. It will use the
// specified sizer function to determine the size of the value. If no sizer
// is specified it will default to a size of one per value. Batches are
// guaranteed to be finished within the specified timeout.
//
// If the source channel is closed the function will send the remaining batch
// and return. If the cancel channel is closed the function will return
// immediately and data may be lost.
func Batch(source <-chan Value, sink chan<- Value, cancel <-chan Signal, sizer func(Value) int, limit int, timeout time.Duration) {
	// prepare slice
	var slice []Value

	// prepare timer
	var timer *time.Timer
	var trigger <-chan time.Time

	// prepare counter
	var counter int

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

			// get value size
			size := 1
			if sizer != nil {
				size = sizer(value)
			}

			// send slice if not enough space
			if counter > 0 && counter+size > limit {
				// send slice
				select {
				case sink <- slice:
				case <-cancel:
				}

				// reset slice
				slice = nil
				counter = 0

				// reset timer if available
				if timer != nil {
					timer.Reset(timeout)
				}
			}

			// add value
			slice = append(slice, value)
			counter += size

			// send slice if already full
			if counter >= limit {
				// send slice
				select {
				case sink <- slice:
				case <-cancel:
				}

				// reset slice
				slice = nil
				counter = 0

				// reset timer if available
				if timer != nil {
					timer.Stop()
					timer = nil
					trigger = nil
				}

				continue
			}

			// set timer if missing
			if timer == nil && timeout > 0 {
				timer = time.NewTimer(timeout)
				trigger = timer.C
			}
		case <-trigger:
			// send slice
			select {
			case sink <- slice:
			case <-cancel:
			}

			// reset slice
			slice = nil
			counter = 0

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
