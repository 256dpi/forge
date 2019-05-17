package forge

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	errTest := errors.New("test")

	reporter := Reporter{}

	var e error
	reporter.Report(func(err error) {
		e = err
	})

	i := 0
	reporter.Repeat(func() error {
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

	reporter := Reporter{}

	assert.PanicsWithValue(t, "missing reporter", func() {
		reporter.Repeat(func() error {
			return errTest
		})
	})
}
