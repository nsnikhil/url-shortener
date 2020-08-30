package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"urlshortner/pkg/shortener/contract"
	"urlshortner/pkg/shortener/service"
)

func ShortenHandler(lgr *zap.Logger, statsd *statsd.Client, svc *service.Service) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		reportAttempt(shortenAPI, statsd)

		var shortenReq contract.ShortenRequest
		err := parseRequest(resp, req, &shortenReq, lgr, shortenAPI, statsd)
		if err != nil {
			return
		}

		shortenResp, err := svc.GetShortenerService().Shorten(shortenReq)
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, false, lgr, shortenAPI, statsd)
			return
		}

		data, err := json.Marshal(&shortenResp)
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, true, lgr, shortenAPI, statsd)
			return
		}

		writeResponse(http.StatusOK, data, resp, lgr)

		reportSuccess(shortenAPI, statsd)
	}
}
