package handler

import (
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"urlshortner/pkg/shortener/service"
)

const locationHeader = "Location"

func RedirectHandler(lgr *zap.Logger, statsd *statsd.Client, svc *service.Service) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		reportAttempt(redirectAPI, statsd)

		longURL, err := svc.GetShortenerService().Redirect(req.URL.Path)
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, false, lgr, redirectAPI, statsd)
			return
		}

		resp.Header().Set(locationHeader, longURL)
		resp.WriteHeader(http.StatusMovedPermanently)

		reportSuccess(redirectAPI, statsd)
	}
}
