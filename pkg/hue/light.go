package hue

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skwair/harmony/optional"
)

// Light is a Hue light bulb.
type Light struct {
	ID               string              `json:"id"`
	State            LightState          `json:"state"`
	SoftWareUpdate   LightSoftwareUpdate `json:"swupdate"`
	Type             string              `json:"type"`
	Name             string              `json:"name"`
	ModelID          string              `json:"modelid"`
	ManufacturerName string              `json:"manufacturername"`
	ProductName      string              `json:"productname"`
	Capabilities     LightCapabilities   `json:"capabilities"`
	Config           LightConfig         `json:"config"`
	UniqueID         string              `json:"uniqueid"`
	SoftWareVersion  string              `json:"swversion"`
	SoftWareConfigID string              `json:"swconfigid"`
	ProductID        string              `json:"productid"`
}

type LightState struct {
	On        bool       `json:"on"`
	Bri       int        `json:"bri"`
	Hue       int        `json:"hue"`
	Sat       int        `json:"sat"`
	Effect    string     `json:"effect"`
	XY        [2]float64 `json:"xy"`
	CT        int        `json:"ct"`
	Alert     string     `json:"alert"`
	ColorMode string     `json:"colormode"`
	Mode      string     `json:"mode"`
	Reachable bool       `json:"reachable"`
}

type LightSoftwareUpdate struct {
	State       string `json:"state"`
	LastInstall string `json:"last_install"`
}

type LightCapabilities struct {
	Certified bool                       `json:"certified"`
	Control   LightCapabilitiesControl   `json:"control"`
	Streaming LightCapabilitiesStreaming `json:"streaming"`
}

type LightCapabilitiesControl struct {
	MinDimLevel    int                        `json:"mindimlevel"`
	MaxLumen       int                        `json:"maxlumen"`
	ColorGamutType string                     `json:"colorgamuttype"`
	ColorGamut     [3][2]float64              `json:"colorgamut"`
	Ct             LightCapabilitiesControlCt `json:"ct"`
}

type LightCapabilitiesControlCt struct {
	Min, Max int
}

type LightCapabilitiesStreaming struct {
	Renderer bool `json:"renderer"`
	Proxy    bool `json:"proxy"`
}

type LightConfig struct {
	ArcheType string `json:"archetype"`
	Function  string `json:"function"`
	Direction string `json:"direction"`
}

type LightConfigStartup struct {
	Mode       string `json:"mode"`
	Configured bool   `json:"configured"`
}

// Lights returns the list of all light bulbs managed by this bridge.
func (c *Client) Lights() ([]Light, error) {
	resp, err := c.doReq(http.MethodGet, "/lights", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]Light
	if err = decode(resp.Body, &res); err != nil {
		return nil, err
	}

	var lights []Light
	for id, l := range res {
		l.ID = id
		lights = append(lights, l)
	}

	return lights, nil
}

// Light returns information about the specified light bulb.
func (c *Client) Light(id string) (*Light, error) {
	endpoint := fmt.Sprintf("/lights/%s", id)
	resp, err := c.doReq(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var light Light
	if err = decode(resp.Body, &light); err != nil {
		return nil, err
	}

	return &light, nil
}

// SetLightStateRequest describes a light state update.
// Only explicitly set fields with be updated.
type SetLightStateRequest struct {
	On             *optional.Bool   `json:"on,omitempty"`
	Bri            *optional.Int    `json:"bri,omitempty"`
	Hue            *optional.Int    `json:"hue,omitempty"`
	Sat            *optional.Int    `json:"sat,omitempty"`
	XY             *[2]float32      `json:"xy,omitempty"`
	CT             *optional.Int    `json:"ct,omitempty"`
	Alert          *optional.String `json:"alert,omitempty"`
	Effect         *optional.String `json:"effect,omitempty"`
	TransitionTime *optional.Int    `json:"transitiontime,omitempty"`
}

// SetLightState sets the state of the specified light bulb.
func (c *Client) SetLightState(id string, req *SetLightStateRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/lights/%s/state", id)
	resp, err := c.doReq(http.MethodPut, endpoint, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeErr(resp.Body)
}

// ToggleLight first queries the state of the specified light bulb then
// switches it on if it was off or off if it was on.
func (c *Client) ToggleLight(id string) error {
	light, err := c.Light(id)
	if err != nil {
		return err
	}

	state := &SetLightStateRequest{
		On: optional.NewBool(!light.State.On),
	}
	return c.SetLightState(id, state)
}
