package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminator(t *testing.T) {
	terminator := Terminator{}

	done1 := make(chan Signal)
	done2 := make(chan Signal)

	go func() {
		<-terminator.Stopping()
		close(done1)
		<-terminator.Killed()
		close(done2)
	}()

	assert.False(t, terminator.IsStopping())
	assert.False(t, terminator.IsKilled())
	assert.NoError(t, terminator.Status())

	terminator.Stop()
	<-done1

	assert.True(t, terminator.IsStopping())
	assert.False(t, terminator.IsKilled())
	assert.Equal(t, ErrStopped, terminator.Status())

	terminator.Kill()
	<-done2

	assert.True(t, terminator.IsStopping())
	assert.True(t, terminator.IsKilled())
	assert.Equal(t, ErrKilled, terminator.Status())
}

func TestTerminatorKill(t *testing.T) {
	terminator := Terminator{}

	done1 := make(chan Signal)
	done2 := make(chan Signal)

	go func() {
		<-terminator.Stopping()
		close(done1)
		<-terminator.Killed()
		close(done2)
	}()

	assert.False(t, terminator.IsStopping())
	assert.False(t, terminator.IsKilled())
	assert.NoError(t, terminator.Status())

	terminator.Kill()
	<-done2

	assert.True(t, terminator.IsStopping())
	assert.True(t, terminator.IsKilled())
	assert.Equal(t, ErrKilled, terminator.Status())
}