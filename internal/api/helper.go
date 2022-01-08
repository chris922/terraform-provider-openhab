package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ReadResponseBody parses the body of the given response as a JSON and unmarshals it into the given target
func ReadResponseBody(resp *http.Response, target interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, target)

	return err
}
