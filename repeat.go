package forge

import "errors"

// TODO: Add repeat with backoff.

// ErrDone is returned to indicate that the task is done.
var ErrDone = errors.New("done")

// Repeat will wrap the provided task and run it repeatedly until ErrDone is
// returned. The specified reporter is called for any other returned error.
func Repeat(task func() error, reporter func(error)) {
	for {
		err := task()
		if err == ErrDone {
			return
		} else if err != nil {
			reporter(err)
		}
	}
}
