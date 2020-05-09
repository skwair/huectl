package hue

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Client is a client that can interact with the API exposed by a Hue bridge.
// Create one with NewClient.
type Client struct {
	url             string
	id              string
	httpClient      *http.Client
	certFingerprint string
}

var defaultHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// NewClient returns a new client for the bridge located at the given IP address
// using the given client ID. It can be further customize with ClientOptions.
// It is strongly recommended to use the HTTPS API (url should start with https://).
func NewClient(url, id string, opts ...ClientOption) *Client {
	c := &Client{
		url:        url,
		id:         id,
		httpClient: defaultHTTPClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ClientOption allows to customize a Hue client.
type ClientOption func(*Client)

// WithHTTPClient overwrites the default HTTP client created to communicate with the bridge.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithCertFingerprint enables TLS certificate fingerprint verification on each request to the bridge.
// Has no effect on non-encrypted connections.
func WithCertFingerprint(fp string) ClientOption {
	return func(c *Client) {
		c.certFingerprint = fp
	}
}
