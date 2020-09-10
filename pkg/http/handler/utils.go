package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"io/ioutil"
	"net/http"
)

const (
	attempt = "attempt"
	success = "success"
	failure = "failure"

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

func handleError(code int, err error, resp http.ResponseWriter, log bool, lgr *zap.Logger, api string, cl *statsd.Client) {
	if log {
		lgr.Error(err.Error())
	}

	reportFailure(api, cl)
	writeResponse(code, bytes.NewBufferString(err.Error()).Bytes(), resp, lgr)
}

func parseRequest(resp http.ResponseWriter, req *http.Request, data interface{}, lgr *zap.Logger, api string, cl *statsd.Client) error {
	if req == nil || req.Body == nil {
		err := errors.New("body is nil")
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, cl)
		return err
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, cl)
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		handleError(http.StatusBadRequest, err, resp, true, lgr, api, cl)
		return err
	}

	return nil
}

func reportAttempt(api string, cl *statsd.Client) {
	incBucket(api, attempt, cl)
}

func reportSuccess(api string, cl *statsd.Client) {
	incBucket(api, success, cl)
}

func reportFailure(api string, cl *statsd.Client) {
	incBucket(api, failure, cl)
}

func incBucket(api, call string, cl *statsd.Client) {
	cl.Increment(fmt.Sprintf("%s.%s.counter", api, call))
}
