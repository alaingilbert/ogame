package ogame

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// OGameClient ...
type OGameClient struct {
	http.Client
	UserAgent    string
	rpsCounter   int32
	rps          int32
	maxRPS       int32
	rpsStartTime int64
}

// NewOGameClient ...
func NewOGameClient() *OGameClient {
	client := &OGameClient{
		Client: http.Client{
			Timeout: 30 * time.Second,
		},
		maxRPS: 0,
	}

	const delay = 1

	go func() {
		for {
			prevRPS := atomic.SwapInt32(&client.rpsCounter, 0)
			atomic.StoreInt32(&client.rps, prevRPS/delay)
			atomic.StoreInt64(&client.rpsStartTime, time.Now().Add(delay*time.Second).UnixNano())
			time.Sleep(delay * time.Second)
		}
	}()

	return client
}

// SetMaxRPS ...
func (c *OGameClient) SetMaxRPS(maxRPS int32) {
	c.maxRPS = maxRPS
}

func (c *OGameClient) incrRPS() {
	newRPS := atomic.AddInt32(&c.rpsCounter, 1)
	if c.maxRPS > 0 && newRPS > c.maxRPS {
		s := atomic.LoadInt64(&c.rpsStartTime) - time.Now().UnixNano()
		// fmt.Printf("throttle %d\n", s)
		time.Sleep(time.Duration(s))
	}
}

// Do executes a request
func (c *OGameClient) Do(req *http.Request) (*http.Response, error) {
	c.incrRPS()
	req.Header.Add("User-Agent", c.UserAgent)
	return c.Client.Do(req)
}

// FakeDo for testing purposes
func (c *OGameClient) FakeDo() {
	c.incrRPS()
	fmt.Println("FakeDo")
}

// GetRPS gets the current client RPS
func (c *OGameClient) GetRPS() int32 {
	return atomic.LoadInt32(&c.rps)
}
