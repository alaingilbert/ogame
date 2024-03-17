package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

// Client special http client that can throttle requests per seconds (RPS).
// Also collect stats about current RPS and overall bytes downloaded/uploaded.
type Client struct {
	sync.Mutex
	*http.Client
	userAgent       string
	rpsCounter      int32 // atomic
	rps             int32 // atomic
	maxRPS          int32 // atomic
	rpsStartTime    int64 // atomic
	bytesDownloaded int64 // atomic
	bytesUploaded   int64 // atomic
}

func (c *Client) BytesDownloaded() int64 {
	return atomic.LoadInt64(&c.bytesDownloaded)
}

func (c *Client) BytesUploaded() int64 {
	return atomic.LoadInt64(&c.bytesUploaded)
}

// NewClient ...
func NewClient(userAgent string) *Client {
	client := &Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRPS:    0,
		userAgent: userAgent,
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
func (c *Client) SetMaxRPS(maxRPS int32) {
	atomic.StoreInt32(&c.maxRPS, maxRPS)
}

func (c *Client) incrRPS() {
	newRPS := atomic.AddInt32(&c.rpsCounter, 1)
	maxRPS := atomic.LoadInt32(&c.maxRPS)
	if maxRPS > 0 && newRPS > maxRPS {
		s := atomic.LoadInt64(&c.rpsStartTime) - time.Now().UnixNano()
		// fmt.Printf("throttle %d\n", s)
		time.Sleep(time.Duration(s))
	}
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.get(url)
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Do executes a request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	c.incrRPS()
	req.Header.Add("User-Agent", c.userAgent)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	atomic.AddInt64(&c.bytesDownloaded, int64(len(body)))
	atomic.AddInt64(&c.bytesUploaded, req.ContentLength)
	// Reset resp.Body so it can be use again
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return resp, err
}

func (c *Client) WithTransport(tr http.RoundTripper, clb func(*Client) error) error {
	c.Lock()
	defer c.Unlock()
	if tr != nil {
		oldTransport := c.Transport
		c.Transport = tr
		defer func() { c.Transport = oldTransport }()
	}
	return clb(c)
}

func (c *Client) SetTransport(tr http.RoundTripper) {
	c.Lock()
	defer c.Unlock()
	c.Transport = tr
}

func (c *Client) UserAgent() string {
	return c.userAgent
}

// FakeDo for testing purposes
func (c *Client) FakeDo() {
	c.incrRPS()
	fmt.Println("FakeDo")
}

// GetRPS gets the current client RPS
func (c *Client) GetRPS() int32 {
	return atomic.LoadInt32(&c.rps)
}
