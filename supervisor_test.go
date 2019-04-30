package forge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupervisor(t *testing.T) {
	s := Supervisor{}

	done1 := make(chan Signal)
	done2 := make(chan Signal)

	go func() {
		<-s.Stopping()
		close(done1)
		<-s.Killed()
		close(done2)
	}()

	assert.False(t, s.IsStopping())
	assert.False(t, s.IsKilled())
	assert.NoError(t, s.Status())

	s.Stop()
	<-done1

	assert.True(t, s.IsStopping())
	assert.False(t, s.IsKilled())
	assert.Equal(t, ErrStopped, s.Status())

	s.Kill()
	<-done2

	assert.True(t, s.IsStopping())
	assert.True(t, s.IsKilled())
	assert.Equal(t, ErrKilled, s.Status())
}

func TestSupervisorKill(t *testing.T) {
	s := Supervisor{}

	done1 := make(chan Signal)
	done2 := make(chan Signal)

	go func() {
		<-s.Stopping()
		close(done1)
		<-s.Killed()
		close(done2)
	}()

	assert.False(t, s.IsStopping())
	assert.False(t, s.IsKilled())
	assert.NoError(t, s.Status())

	s.Kill()
	<-done2

	assert.True(t, s.IsStopping())
	assert.True(t, s.IsKilled())
	assert.Equal(t, ErrKilled, s.Status())
}
