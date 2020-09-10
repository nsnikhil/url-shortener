package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"urlshortner/pkg/reporters"
)

const (
	pingAPI     = "ping"
	shortenAPI  = "shorten"
	redirectAPI = "redirect"
)

func writeResponse(code int, data []byte, resp http.ResponseWriter, lgr *zap.Logger) {
	resp.WriteHeader(code)
	if _, err := resp.Write(data); err != nil {
		lgr.Error(err.Error())
	}
}

func handleError(code int, err error, resp http.ResponseWriter, log bool, lgr *zap.Logger, api string, sc reporters.StatsDClient) {
	if log {
		lgr.Error(err.Error())
	}

	sc.ReportFailure(api)
	writeResponse(code, bytes.NewBufferString(err.Error()).Bytes(), resp, lgr)
}

func parseRequest(resp http.ResponseWriter, req *http.Request, data interface{}, lgr *zap.Logger, api string, sc reporters.StatsDClient) error {
	if req == nil || req.Body == nil {
		err := errors.New("body is nil")
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, sc)
		return err
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, sc)
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, sc)
		return err
	}

	return nil
}
