package forge

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	in := make(chan Value)
	out := make(chan Value)

	go func() {
		Batch(in, out, nil, 3, 0)
		close(out)
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
		}

		close(in)
	}()

	var list []Value
	for i := range out {
		list = append(list, i)
	}

	assert.Equal(t, []Value{[]Value{1, 2, 3}, []Value{4, 5, 6}, []Value{7, 8, 9}, []Value{10}}, list)
}

func TestBatchTimeout(t *testing.T) {
	in := make(chan Value)
	out := make(chan Value)

	go func() {
		Batch(in, out, nil, 3, 3*time.Millisecond)
		close(out)
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
			time.Sleep(2 * time.Millisecond)
		}

		close(in)
	}()

	var list []Value
	for i := range out {
		list = append(list, i)
	}

	assert.Equal(t, []Value{[]Value{1, 2}, []Value{3, 4}, []Value{5, 6}, []Value{7, 8}, []Value{9, 10}}, list)
}

func TestBatchCancel(t *testing.T) {
	in := make(chan Value)
	out := make(chan Value)
	cancel := make(chan Signal)

	go func() {
		Batch(in, out, cancel, 3, 0)
		close(out)
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
		}

		close(cancel)
	}()

	var list []Value
	for i := range out {
		list = append(list, i)
	}

	assert.Equal(t, []Value{[]Value{1, 2, 3}, []Value{4, 5, 6}, []Value{7, 8, 9}}, list)
}

func BenchmarkBatch(b *testing.B) {
	size := 1000

	in := make(chan Value, size)
	out := make(chan Value, b.N/size+1)

	go Batch(in, out, nil, size, 0)

	for i := 0; i < b.N; i++ {
		in <- i
	}

	close(in)
}

func BenchmarkBatchTimeout(b *testing.B) {
	size := 1000

	in := make(chan Value, size)
	out := make(chan Value, b.N/size+1)

	go Batch(in, out, nil, size, time.Second)

	for i := 0; i < b.N; i++ {
		in <- i
	}

	close(in)
}
