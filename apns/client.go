// Package apns is a go Apple Push Notification Service (APNs) provider that
// allows you to send remote notifications to your iOS, tvOS, and OS X
// apps, using the new APNs HTTP/2 network protocol.
package apns

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
)

const (
	// ApnsPriorityHigh Apns Priority High
	ApnsPriorityHigh = 10
	// ApnsPriorityLow Apns Priority Low
	ApnsPriorityLow = 5
)

const (
	// StatusCodesSuccess Success
	StatusCodesSuccess = 200
	// StatusCodesBadReq Bad request
	StatusCodesBadReq = 400
	// StatusCodesCerErr There was an error with the certificate.
	StatusCodesCerErr = 403
	// StatusCodesMethodErr The request used a bad :method value. Only POST requests are supported.
	StatusCodesMethodErr = 405
	// StatusCodesNoActive The device token is no longer active for the topic.
	StatusCodesNoActive = 410
	// StatusCodesPayloadTooLarge The notification payload was too large.
	StatusCodesPayloadTooLarge = 413
	// StatusCodesTooManyReq The server received too many requests for the same device token.
	StatusCodesTooManyReq = 429
	// StatusCodesServerErr Internal server error
	StatusCodesServerErr = 500
	// StatusCodesServerUnavailable The server is shutting down and unavailable.
	StatusCodesServerUnavailable = 503
)

// Apple HTTP/2 Development & Production urls
const (
	HostDevelopment = "https://api.development.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

// DefaultHost is a mutable var for testing purposes
var DefaultHost = HostDevelopment

// Client represents a connection with the APNs
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
	BoundID     string
}

// NewClient returns a new Client with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
func NewClient(certificate tls.Certificate, timeout time.Duration) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &Client{
		HTTPClient:  &http.Client{Transport: transport, Timeout: timeout},
		Certificate: certificate,
		Host:        DefaultHost,
	}
}

// NewClientWithProxy new client with proxy.
func NewClientWithProxy(certificate tls.Certificate, timeout, deadlineTime time.Duration, proxyURL string) *Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	return &Client{
		HTTPClient:  &http.Client{Transport: proxyTransport(proxyURL, tlsConfig, timeout, deadlineTime), Timeout: timeout},
		Certificate: certificate,
		Host:        DefaultHost,
	}
}

func proxyTransport(proxyAddr string, config *tls.Config, timeout, deadlineTime time.Duration) *http2.Transport {
	return &http2.Transport{
		DialTLS: func(network, addr string, cfg *tls.Config) (conn net.Conn, err error) {
			dialer := &net.Dialer{Timeout: timeout / 2}
			d, err := proxy.SOCKS5("tcp", proxyAddr, nil, dialer)
			if err != nil {
				return nil, err
			}
			conn, err = d.Dial(network, addr)
			if err != nil {
				fmt.Println(err)
				conn, err = dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
			}
			if err = conn.SetDeadline(time.Now().Add(deadlineTime)); err != nil {
				return nil, err
			}
			cn := tls.Client(conn, cfg)
			// from dialTLSDefault
			if err := cn.Handshake(); err != nil {
				return nil, err
			}
			if !cfg.InsecureSkipVerify {
				if err := cn.VerifyHostname(cfg.ServerName); err != nil {
					return nil, err
				}
			}
			state := cn.ConnectionState()
			if p := state.NegotiatedProtocol; p != http2.NextProtoTLS {
				return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, http2.NextProtoTLS)
			}
			if !state.NegotiatedProtocolIsMutual {
				return nil, errors.New("http2: could not negotiate protocol mutually")
			}
			return cn, nil
		},
		TLSClientConfig: config,
	}
}

// Development sets the Client to use the APNs development push endpoint.
func (c *Client) Development() *Client {
	c.Host = HostDevelopment
	return c
}

// Production sets the Client to use the APNs production push endpoint.
func (c *Client) Production() *Client {
	c.Host = HostProduction
	return c
}

// Push sends a Notification to the APNs gateway. If the underlying http.Client
// is not currently connected, this method will attempt to reconnect
// transparently before sending the notification.
func (c *Client) Push(deviceToken string, payload *Payload, overTime int64) (response *Response, err error) {
	var (
		url     string
		req     *http.Request
		httpRes *http.Response
	)
	url = fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(payload.Marshal()))
	req.Header.Set("apns-topic", c.BoundID)
	req.Header.Set("apns-expiration", strconv.FormatInt(overTime, 10))
	httpRes, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()
	response = &Response{}
	response.StatusCode = httpRes.StatusCode
	response.ApnsID = httpRes.Header.Get("apns-id")
	return response, nil
}
