package handler

import (
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/reporters"
)

const locationHeader = "Location"

func RedirectHandler(lgr *zap.Logger, statsdClient reporters.StatsDClient, elongator elongator.Elongator) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		statsdClient.ReportAttempt(redirectAPI)

		longURL, err := elongator.Elongate(req.URL.Path[1:])
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, false, lgr, redirectAPI, statsdClient)
			return
		}

		resp.Header().Set(locationHeader, longURL)
		resp.WriteHeader(http.StatusMovedPermanently)

		statsdClient.ReportSuccess(redirectAPI)
	}
}
