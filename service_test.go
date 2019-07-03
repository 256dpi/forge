package forge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("tes")

func TestService(t *testing.T) {
	service := Service{}

	var errs []error
	service.Report(func(err error) {
		errs = append(errs, err)
	})

	i := 0
	service.Run(1, func() error {
		i++
		if i == 3 {
			return ErrDone
		}

		return testErr
	}, service.Stop)

	<-service.Stopping()
	<-service.Done()

	service.Kill()
	<-service.Killed()

	assert.Equal(t, 3, i)
	assert.Equal(t, []error{testErr, testErr}, errs)
}
