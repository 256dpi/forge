package forge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	service := Service{}

	var err error
	service.Report(func(e error) {
		err = e
	})

	i := 0
	service.Run(1, func() error {
		i++
		if i == 2 {
			return ErrDone
		}
		return errors.New("foo")
	}, service.Stop)

	<-service.Stopping()
	<-service.Done()

	service.Kill()
	<-service.Killed()

	assert.Equal(t, 2, i)
	assert.Error(t, err)
}
