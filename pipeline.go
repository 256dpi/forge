package forge

import "time"

// Pipeline extends manager with data pipeline primitives.
type Pipeline struct {
	Manager

	input  chan Value
	source chan Value
}

// Open will open the pipeline and allow external input.
func (p *Pipeline) Open(buffer int) {
	if p.source != nil {
		panic("only one source allowed")
	}

	p.input = make(chan Value, buffer)
	p.source = p.input
}

// Input will return the external input channel.
func (p *Pipeline) Input() chan<- Value {
	return p.input
}

// Source will run the source task that fills the provided channel with values.
func (p *Pipeline) Source(n int, fn func(chan<- Value), buffer int) {
	if p.source != nil {
		panic("only one source allowed")
	}

	out := make(chan Value, buffer)

	p.Run(n, func() {
		fn(out)
	}, func() {
		close(out)
	})

	p.source = out
}

// Filter is an intermediary task that processes values.
func (p *Pipeline) Filter(n int, fn func(<-chan Value, chan<- Value), buffer int) {
	if p.source == nil {
		panic("missing source")
	}

	in := p.source
	out := make(chan Value, buffer)

	p.Run(n, func() {
		fn(in, out)
	}, func() {
		close(out)
	})

	p.source = out
}

// FilterFunc augments Filter by running the specified function once per
// received value.
func (p *Pipeline) FilterFunc(n int, fn func(Value, chan<- Value), buffer int) {
	p.Filter(n, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				// call function
				fn(v, out)
			}
		}
	}, buffer)
}

// Batch is an intermediary task that batches values up.
func (p *Pipeline) Batch(n int, sizer func(Value) int, limit int, stop <-chan Signal, timeout time.Duration, buffer int) {
	if p.source == nil {
		panic("missing source")
	}

	in := p.source
	out := make(chan Value, buffer)

	p.Run(n, func() {
		Batch(in, out, stop, sizer, limit, timeout)
	}, func() {
		close(out)
	})

	p.source = out
}

// Sink is the final task that receives all processed values.
func (p *Pipeline) Sink(n int, fn func(<-chan Value)) {
	if p.source == nil {
		panic("missing source")
	}

	in := p.source

	p.Run(n, func() {
		fn(in)
	}, nil)

	p.source = nil
}

// SinkFunc augments Sink by running the specified function once per received
// value.
func (p *Pipeline) SinkFunc(n int, fn func(Value)) {
	p.Sink(n, func(in <-chan Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				// call function
				fn(v)
			}
		}
	})
}

// Output will return the output channel.
func (p *Pipeline) Output() <-chan Value {
	if p.source == nil {
		panic("missing source")
	}

	return p.source
}
