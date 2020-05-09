package hue

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) doReq(method, endpoint string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/%s%s", c.url, c.id, endpoint)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.TLS != nil && c.certFingerprint != "" {
		fp := computeFingerprint(resp.TLS.PeerCertificates[0].Raw)
		if c.certFingerprint != fp {
			return nil, errors.New("certificate fingerprint mismatch")
		}
	}

	return resp, nil
}
