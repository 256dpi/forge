package forge

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	done := make(chan struct{})

	var counter int64

	Run(3, func() {
		atomic.AddInt64(&counter, 1)
	}, func() {
		close(done)
	})

	<-done

	assert.Equal(t, int64(3), counter)
}
