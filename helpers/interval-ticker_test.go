package helpers

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetInterval(t *testing.T) {
	type TimesCalled struct {
		sync.RWMutex
		value int
	}
	timesCalled := TimesCalled{
		value: 0,
	}
	timer := time.NewTimer(time.Duration(1100) * time.Millisecond)
	ticker := SetInterval(func() {
		timesCalled.Lock()
		timesCalled.value++
		timesCalled.Unlock()
	}, 250)
	<-timer.C

	ClearInterval(ticker)

	assert.Equal(t, 4, timesCalled.value, "should call function over duration of interval")
}
