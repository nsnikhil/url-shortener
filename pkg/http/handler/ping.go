package handler

import (
	"bytes"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
)

func PingHandler(lgr *zap.Logger, statsd *statsd.Client) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		reportAttempt(pingAPI, statsd)

		writeResponse(http.StatusOK, bytes.NewBufferString("pong").Bytes(), resp, lgr)

		reportSuccess(pingAPI, statsd)
	}
}
