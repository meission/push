package mipush

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/meission/log4go"
)

type Client struct {
	Header     http.Header
	HTTPClient *http.Client
	Url        string
}

// NewClient:xm push client
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

// Production:  Production URL
func (c *Client) SetProductionUrl(url string) {
	c.Url = HostProBase + url
}

// Development:  Production URL
func (c *Client) SetDevelopmentUrl(url string) {
	c.Url = HostDevbase + url
}

// push
func (c *Client) Push(xm *XMMessage) (response *Response, err error) {
	var (
		req  *http.Request
		resp *http.Response
		body []byte
	)
	req, err = http.NewRequest("POST", c.Url, bytes.NewBuffer([]byte(xm.xmuv.Encode())))
	if err != nil {
		log.Error("http.NewRequest(\"POST\",c.Url:%s}error:%v", c.Url, err)
		return nil, err
	}
	req.Header = c.Header
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		log.Error("c.HTTPClient.Do(),error:%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("ioutil.ReadAll(resp.Body);is error(%v).", err)
		return nil, err
	}
	response = &Response{}
	err = response.Unmarshal(body)
	if err != nil {
		log.Error("response.Unmarshal(body),error:%v.", err)
		return nil, err
	}
	return response, nil
}
