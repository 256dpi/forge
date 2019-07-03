package forge

import (
	"errors"
	"fmt"
)

func Example() {
	// prepare service
	service := &Service{}

	// set reporter
	var errs []error
	service.Report(func(err error) {
		errs = append(errs, err)
	})

	// run task
	i := 0
	service.Run(1, func() error {
		i++

		if i == 5 {
			return ErrDone
		}
		if i%2 == 0 {
			return errors.New("foo")
		}

		return nil
	}, func() {
		fmt.Println("finalize")
	})

	// wait for exit
	<-service.Done()

	// print output
	fmt.Println(i)
	fmt.Println(errs)

	// Output:
	// finalize
	// 5
	// [foo foo]
}
