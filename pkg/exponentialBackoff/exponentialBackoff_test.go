package exponentialBackoff

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
		clock.BlockUntil(1)
		clock.Advance(8000 * time.Millisecond)
		clock.BlockUntil(0)
		atomic.AddUint32(&counter, 1)
		wg.Done()
	}()
	e := New(context.Background(), clock, 60)
	e.Wait() // Wait 1s
	e.Wait() // Wait 2s
	e.Wait() // Wait 4s
	e.Wait() // Wait 8s
	wg.Wait()
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}

//func TestExponentialBackoff_1(t *testing.T) {
//	start := time.Now()
//	clock := clockwork.NewRealClock()
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	e := New(ctx, clock, 60)
//	go func() {
//		time.Sleep(5 * time.Second)
//		e.Reset()
//	}()
//	for range e.Iterator() {
//		fmt.Println("?????", time.Since(start))
//	}
//	assert.Equal(t, 1, 2)
//}
