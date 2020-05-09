package hue

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const hueDiscoveryURL = "https://discovery.meethue.com"

// Bridge is a Hue bridge, discovered on a local network by DiscoverBridges.
type Bridge struct {
	ID              string
	IPAddr          string
	Name            string
	CertFingerprint string
}

// DiscoverBridges searches for Hue bridges on the local network using Philips'
// Hue service discovery system. Note that the given HTTP client is only used
// to connect to Hue bridges, not to contact the discovery endpoint.
func DiscoverBridges(httpClient *http.Client) ([]Bridge, error) {
	c := http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(hueDiscoveryURL)
	if err != nil {
		return nil, fmt.Errorf("unable to send discovery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery request: non-200 response: %s", http.StatusText(resp.StatusCode))
	}

	var bridgeInfos []struct {
		ID     string `json:"id"`
		IPAddr string `json:"internalipaddress"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&bridgeInfos); err != nil {
		return nil, fmt.Errorf("unable to decode discovery response: %w", err)
	}

	var bridges []Bridge
	for _, info := range bridgeInfos {
		name, fp, err := pingBridge(httpClient, info.IPAddr)
		if err != nil {
			fmt.Printf("unable to ping bridge %q at %s: %v, skipping it\n", info.ID, info.IPAddr, err)
			continue
		}

		bridges = append(bridges, Bridge{
			ID:              info.ID,
			IPAddr:          info.IPAddr,
			Name:            name,
			CertFingerprint: fp,
		})
	}

	return bridges, nil
}

func pingBridge(httpClient *http.Client, ip string) (name, fingerprint string, err error) {
	resp, err := httpClient.Get(fmt.Sprintf("https://%s/api/config", ip))
	if err != nil {
		return "", "", fmt.Errorf("unable to get Hue bridge information: %w", err)
	}
	defer resp.Body.Close()

	var b struct {
		Name string `json:"name"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return "", "", fmt.Errorf("unable to decode Hue bridge information: %w", err)
	}

	if len(resp.TLS.PeerCertificates) != 1 {
		return "", "", fmt.Errorf("expected exactly one peer certificate; got %d", len(resp.TLS.PeerCertificates))
	}

	return b.Name, computeFingerprint(resp.TLS.PeerCertificates[0].Raw), nil
}

func computeFingerprint(d []byte) string {
	sum := sha1.Sum(d)

	var s strings.Builder
	for i, byt := range sum {
		s.WriteString(hex.EncodeToString([]byte{byt}))

		if i+1 < len(sum) {
			s.WriteByte(':')
		}
	}

	return s.String()
}
