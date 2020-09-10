package handler

import (
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"urlshortner/pkg/elongator"
)

const locationHeader = "Location"

func RedirectHandler(lgr *zap.Logger, statsd *statsd.Client, elongator elongator.Elongator) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		reportAttempt(redirectAPI, statsd)

		longURL, err := elongator.Elongate(req.URL.Path[1:])
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, false, lgr, redirectAPI, statsd)
			return
		}

		resp.Header().Set(locationHeader, longURL)
		resp.WriteHeader(http.StatusMovedPermanently)

		reportSuccess(redirectAPI, statsd)
	}
}
