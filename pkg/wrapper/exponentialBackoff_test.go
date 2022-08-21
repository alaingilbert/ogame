package wrapper

import (
	"context"
	"github.com/alaingilbert/clockwork"
	"github.com/magiconair/properties/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExponentialBackoff_Wait(t *testing.T) {
	var counter uint32
	clock := clockwork.NewFakeClock()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		clock.BlockUntil(1)
		clock.Advance(1000 * time.Millisecond)
		clock.BlockUntil(0)
		clock.BlockUntil(1)
		clock.Advance(2000 * time.Millisecond)
		clock.BlockUntil(0)
		clock.BlockUntil(1)
		clock.Advance(4000 * time.Millisecond)
		clock.BlockUntil(0)
		atomic.AddUint32(&counter, 1)
	}()
	go func() {
		e := NewExponentialBackoff(context.Background(), 60)
		e.SetClock(clock)
		e.Wait() // First time has no wait
		e.Wait() // Wait 1s
		e.Wait() // Wait 2s
		e.Wait() // Wait 4s
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}
