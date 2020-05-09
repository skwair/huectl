package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RegisterUser registers a new user on the Hue bridge at the given address.
// The bridge needs to be in "link" mode, i.e. the button must have been pressed
// in the last 30 seconds for the registration to work.
// The device type is an identifier for this new user and must have the following
// format: <application_name>#<device_name>. See https://developers.meethue.com/develop/hue-api/7-configuration-api/#create-user
// for more information.
func RegisterUser(httpClient *http.Client, addr, deviceType string) (username string, err error) {
	registerReq := struct {
		DeviceType string `json:"devicetype"`
	}{
		DeviceType: deviceType,
	}
	b, err := json.Marshal(registerReq)
	if err != nil {
		return "", fmt.Errorf("unable to marshal register user request body: %w", err)
	}

	resp, err := httpClient.Post(fmt.Sprintf("https://%s/api", addr), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("unable to send second register user request: %w", err)
	}
	defer resp.Body.Close()

	registerResp := make([]struct {
		Success struct {
			Username string `json:"username"`
		} `json:"success"`
	}, 1)
	if err = json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return "", fmt.Errorf("unable to decode register response: %w", err)
	}

	if len(registerResp) == 0 || registerResp[0].Success.Username == "" {
		return "", errors.New("failed to register new user")
	}

	return registerResp[0].Success.Username, nil
}
