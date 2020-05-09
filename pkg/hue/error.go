package hue

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

// ErrorSet is a set of API errors returned by the Hue API.
type ErrorSet []Error

// Error implements the `error` interface.
func (s ErrorSet) Error() string {
	var sb strings.Builder

	for _, err := range s {
		sb.WriteString(err.Error())
	}

	return sb.String()
}

type errResp struct {
	Error Error `json:"error"`
}

// Error is an API error returned by the Hue API.
type Error struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Error implements the `error` interface.
func (e Error) Error() string { return e.Description }

func decode(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if bytes.HasPrefix(data, []byte(`[{"error":{`)) {
		var errs []errResp
		if err = json.Unmarshal(data, &errs); err != nil {
			return err
		}

		var es ErrorSet
		for _, e := range errs {
			es = append(es, e.Error)
		}

		return es
	}

	return json.Unmarshal(data, v)
}

func decodeErr(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if bytes.HasPrefix(data, []byte(`[{"error":{`)) {
		var errs []errResp
		if err = json.Unmarshal(data, &errs); err != nil {
			return err
		}

		var es ErrorSet
		for _, e := range errs {
			es = append(es, e.Error)
		}

		return es
	}

	return nil
}
