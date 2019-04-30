package forge

// A Reporter manages repeating tasks and their error reporting.
type Reporter struct {
	fn func(error)
}

// Report will set the reporter function.
func (r *Reporter) Report(fn func(error)) {
	r.fn = fn
}

// Repeat will repeat the provided function.
func (r *Reporter) Repeat(fn func() error) {
	// check reporter
	if r.fn == nil {
		panic("missing reporter")
	}

	Repeat(fn, r.fn)
}
