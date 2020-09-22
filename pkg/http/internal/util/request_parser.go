package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func ParseRequest(req *http.Request, data interface{}) error {
	if req == nil {
		return errors.New("request is nil")
	}

	if req.Body == nil {
		return errors.New("request body is nil")
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}
