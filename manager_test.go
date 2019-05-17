package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	manager := Manager{}

	start := 0
	stop := 0

	manager.Run(3, func() {
		start++
	}, func() {
		stop++
	})

	assert.False(t, manager.Finished())
	<-manager.Done()
	assert.True(t, manager.Finished())

	manager.Run(3, func() {
		start++
	}, func() {
		stop++
	})

	assert.False(t, manager.Finished())
	<-manager.Done()
	assert.True(t, manager.Finished())

	assert.Equal(t, 6, start)
	assert.Equal(t, 2, stop)
}
