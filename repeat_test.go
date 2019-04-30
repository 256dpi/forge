package forge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	errTest := errors.New("test")

	i1 := 0
	var e1 error
	Repeat(func() error {
		i1++
		if i1 == 2 {
			return ErrDone
		}
		return nil
	}, func(e error) {
		e1 = e
	})
	assert.Equal(t, 2, i1)
	assert.NoError(t, e1)

	i2 := 0
	var e2 error
	Repeat(func() error {
		i2++
		if i2 == 2 {
			return ErrDone
		}
		return errTest
	}, func(e error) {
		e2 = e
	})
	assert.Equal(t, 2, i2)
	assert.Equal(t, errTest, e2)
}
