package mipush

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Header     http.Header
	HTTPClient *http.Client
	URL        string
}

// NewClient xm push client
// set request header and auth
func NewClient(auth string) *Client {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Authorization", AUTH_PREFIX+auth)
	return &Client{
		Header:     header,
		HTTPClient: &http.Client{},
	}
}

// SetProductionURL   Production URL
func (c *Client) SetProductionURL(url string) {
	c.Url = HostProBase + url
}

// SetDevelopmentURL   Production URL .
func (c *Client) SetDevelopmentURL(url string) {
	c.Url = HostDevbase + url
}

// Push xm message .
func (c *Client) Push(xm *XMMessage) (response *Response, err error) {
	var (
		req  *http.Request
		resp *http.Response
		body []byte
	)
	req, err = http.NewRequest("POST", c.URL, bytes.NewBuffer([]byte(xm.xmuv.Encode())))
	if err != nil {
		return nil, err
	}
	req.Header = c.Header
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	response = &Response{}
	err = response.Unmarshal(body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
