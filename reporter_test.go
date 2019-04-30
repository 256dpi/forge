package forge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	errTest := errors.New("test")

	r := Reporter{}

	var e error
	r.Report(func(err error) {
		e = err
	})

	i := 0
	r.Repeat(func() error {
		i++
		if i == 2 {
			return ErrDone
		}
		return errTest
	})

	assert.Equal(t, 2, i)
	assert.Equal(t, errTest, e)
}

func TestReporterMissing(t *testing.T) {
	errTest := errors.New("test")

	r := Reporter{}

	assert.PanicsWithValue(t, "missing reporter", func() {
		r.Repeat(func() error {
			return errTest
		})
	})
}
