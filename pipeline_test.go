package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	pipeline := Pipeline{}

	pipeline.Source(1, func(out chan<- Value) {
		i := 0

		for {
			i++

			out <- 21

			if i == 6 {
				return
			}
		}
	}, 0)

	pipeline.Filter(3, func(in <-chan Value, out chan<- Value) {
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

	pipeline.Batch(1, nil, 2, nil, 0, 0)

	var list [][]Value

	pipeline.Sink(1, func(in <-chan Value) {
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

	<-pipeline.Done()

	assert.Equal(t, [][]Value{{42, 42}, {42, 42}, {42, 42}}, list)
}

func TestPipelineOpen(t *testing.T) {
	pipeline := Pipeline{}

	pipeline.Open(3)

	pipeline.Filter(1, func(in <-chan Value, out chan<- Value) {
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

	pipeline.Input() <- 1
	pipeline.Input() <- 2
	pipeline.Input() <- 3
	close(pipeline.Input())

	v1 := <-pipeline.Output()
	v2 := <-pipeline.Output()
	v3 := <-pipeline.Output()
	<-pipeline.Done()

	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
}

func TestPipelineFunc(t *testing.T) {
	pipeline := Pipeline{}

	pipeline.Open(3)

	pipeline.FilterFunc(1, func(v Value, out chan<- Value) {
		out <- v.(int) * 2
	}, 0)

	var list []int
	pipeline.SinkFunc(1, func(v Value) {
		list = append(list, v.(int))
	})

	pipeline.Input() <- 1
	pipeline.Input() <- 2
	pipeline.Input() <- 3
	close(pipeline.Input())

	<-pipeline.Done()

	assert.Equal(t, []int{2, 4, 6}, list)
}

func BenchmarkPipelineSingle(b *testing.B) {
	pipeline := Pipeline{}

	pipeline.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 0)

	pipeline.Filter(1, func(in <-chan Value, out chan<- Value) {
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

	pipeline.Batch(1, nil, 2, nil, 0, 0)

	pipeline.Sink(1, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-pipeline.Done()
}

func BenchmarkPipelineDistributed(b *testing.B) {
	pipeline := Pipeline{}

	pipeline.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 0)

	pipeline.Filter(10, func(in <-chan Value, out chan<- Value) {
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

	pipeline.Batch(10, nil, 2, nil, 0, 0)

	pipeline.Sink(10, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-pipeline.Done()
}

func BenchmarkPipelineBuffered(b *testing.B) {
	pipeline := Pipeline{}

	pipeline.Source(1, func(out chan<- Value) {
		i := 0
		for {
			i++
			out <- 21
			if i == b.N {
				return
			}
		}
	}, 1000)

	pipeline.Filter(3, func(in <-chan Value, out chan<- Value) {
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

	pipeline.Batch(1, nil, 2, nil, 0, 1000)

	pipeline.Sink(1, func(in <-chan Value) {
		for {
			select {
			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	})

	<-pipeline.Done()
}
