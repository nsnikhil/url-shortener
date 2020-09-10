package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/reporters"
	"urlshortner/pkg/shortener"
)

func ShortenHandler(lgr *zap.Logger, statsdClient reporters.StatsDClient, shortener shortener.Shortener) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		statsdClient.ReportAttempt(shortenAPI)

		var shortenReq contract.ShortenRequest
		err := parseRequest(resp, req, &shortenReq, lgr, shortenAPI, statsdClient)
		if err != nil {
			return
		}

		shortURL, err := shortener.Shorten(shortenReq.URL)
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, false, lgr, shortenAPI, statsdClient)
			return
		}

		shortenResp := contract.ShortenResponse{ShortURL: shortURL}

		data, err := json.Marshal(&shortenResp)
		if err != nil {
			handleError(http.StatusInternalServerError, err, resp, true, lgr, shortenAPI, statsdClient)
			return
		}

		writeResponse(http.StatusOK, data, resp, lgr)

		statsdClient.ReportSuccess(shortenAPI)
	}
}
