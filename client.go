package ogame

import "net/http"

type ogameClient struct {
	http.Client
	UserAgent string
}

func (c *ogameClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", c.UserAgent)
	return c.Client.Do(req)
}
