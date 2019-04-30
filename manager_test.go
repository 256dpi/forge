package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	s := Manager{}

	start := 0
	stop := 0

	s.Run(3, func() {
		start++
	}, func() {
		stop++
	})

	assert.False(t, s.Finished())
	<-s.Done()
	assert.True(t, s.Finished())

	s.Run(3, func() {
		start++
	}, func() {
		stop++
	})

	assert.False(t, s.Finished())
	<-s.Done()
	assert.True(t, s.Finished())

	assert.Equal(t, 6, start)
	assert.Equal(t, 2, stop)
}
