package handler

import (
	"bytes"
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/reporters"
)

func PingHandler(lgr *zap.Logger, statsdClient reporters.StatsDClient) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		statsdClient.ReportAttempt(pingAPI)

		writeResponse(http.StatusOK, bytes.NewBufferString("pong").Bytes(), resp, lgr)

		statsdClient.ReportSuccess(pingAPI)
	}
}
