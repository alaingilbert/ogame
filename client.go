package ogame

import "net/http"

// OGameClient ...
type OGameClient struct {
	http.Client
	UserAgent string
}

// Do ...
func (c *OGameClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", c.UserAgent)
	return c.Client.Do(req)
}
