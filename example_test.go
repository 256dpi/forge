package forge

import (
	"fmt"
	"time"
)

func Example() {
	service := &Service{}

	service.Report(func(err error) {
		panic(err.Error())
	})

	i := 0
	service.Source(1, func(values chan<- Value) error {
		i++
		values <- i
		if i == 5 {
			return ErrDone
		}
		return nil
	}, 1)

	service.Batch(1, nil, 3, time.Millisecond, 1)

	service.FilterFunc(1, func(v Value, out chan<- Value) error {
		out <- v
		return nil
	}, 1)

	var out []Value
	service.SinkFunc(1, func(v Value) error {
		out = append(out, v)
		return nil
	})

	<-service.Done()

	fmt.Printf("%+v\n", out)

	// Output:
	// [[1 2 3] [4 5]]
}
