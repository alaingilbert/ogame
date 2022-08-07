package ogame

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// OGameClient ...
type OGameClient struct {
	*http.Client
	UserAgent       string
	rpsCounter      int32 // atomic
	rps             int32 // atomic
	maxRPS          int32 // atomic
	rpsStartTime    int64 // atomic
	bytesDownloaded int64
	bytesUploaded   int64
}

// NewOGameClient ...
func NewOGameClient() *OGameClient {
	client := &OGameClient{
		Client: &http.Client{
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
	atomic.StoreInt32(&c.maxRPS, maxRPS)
}

func (c *OGameClient) incrRPS() {
	newRPS := atomic.AddInt32(&c.rpsCounter, 1)
	maxRPS := atomic.LoadInt32(&c.maxRPS)
	if maxRPS > 0 && newRPS > maxRPS {
		s := atomic.LoadInt64(&c.rpsStartTime) - time.Now().UnixNano()
		// fmt.Printf("throttle %d\n", s)
		time.Sleep(time.Duration(s))
	}
}

func (c *OGameClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Do executes a request
func (c *OGameClient) Do(req *http.Request) (*http.Response, error) {
	c.incrRPS()
	req.Header.Add("User-Agent", c.UserAgent)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	c.bytesDownloaded += int64(len(body))
	c.bytesUploaded += req.ContentLength
	// Reset resp.Body so it can be use again
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return resp, err
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
