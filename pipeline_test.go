package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	p := Pipeline{}

	p.Source(1, func(out chan<- Value) {
		i := 0

		for {
			i++

			out <- 21

			if i == 6 {
				return
			}
		}
	}, 0)

	p.Filter(3, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v.(int) * 2
			}
		}
	}, 0)

	p.Batch(1, 2, nil, 0, 0)

	var list [][]Value

	p.Sink(1, func(in <-chan Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				list = append(list, v.([]Value))
			}
		}
	})

	<-p.Done()

	assert.Equal(t, [][]Value{{42, 42}, {42, 42}, {42, 42}}, list)
}

func TestPipelineOpen(t *testing.T) {
	p := Pipeline{}

	p.Open(3)

	p.Filter(1, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v.(int) * 2
			}
		}
	}, 0)

	p.Input() <- 1
	p.Input() <- 2
	p.Input() <- 3
	close(p.Input())

	v1 := <-p.Output()
	v2 := <-p.Output()
	v3 := <-p.Output()
	<-p.Done()

	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
}

func TestPipelineFunc(t *testing.T) {
	p := Pipeline{}

	p.Open(3)

	p.FilterFunc(1, func(v Value, out chan<- Value) {
		out <- v.(int) * 2
	}, 0)

	var list []int
	p.SinkFunc(1, func(v Value) {
		list = append(list, v.(int))
	})

	p.Input() <- 1
	p.Input() <- 2
	p.Input() <- 3
	close(p.Input())

	<-p.Done()

	assert.Equal(t, []int{2, 4, 6}, list)
}

func BenchmarkPipelineSingle(b *testing.B) {
	p := Pipeline{}

	p.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 0)

	p.Filter(1, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v.(int) * 2
			}
		}
	}, 0)

	p.Batch(1, 2, nil, 0, 0)

	p.Sink(1, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-p.Done()
}

func BenchmarkPipelineDistributed(b *testing.B) {
	p := Pipeline{}

	p.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 0)

	p.Filter(10, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v.(int) * 2
			}
		}
	}, 0)

	p.Batch(10, 2, nil, 0, 0)

	p.Sink(10, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-p.Done()
}

func BenchmarkPipelineBuffered(b *testing.B) {
	p := Pipeline{}

	p.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 1000)

	p.Filter(3, func(in <-chan Value, out chan<- Value) {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v.(int) * 2
			}
		}
	}, 1000)

	p.Batch(1, 2, nil, 0, 1000)

	p.Sink(1, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-p.Done()
}
