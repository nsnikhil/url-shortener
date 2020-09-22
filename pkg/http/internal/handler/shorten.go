package handler

import (
	"net/http"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/http/internal/liberr"
	"urlshortner/pkg/http/internal/util"
	"urlshortner/pkg/shortener"
)

type ShortenHandler struct {
	shortener shortener.Shortener
}

func (sh *ShortenHandler) Shorten(resp http.ResponseWriter, req *http.Request) error {
	var shortenReq contract.ShortenRequest
	err := util.ParseRequest(req, &shortenReq)
	if err != nil {
		return liberr.ValidationError(err.Error())
	}

	shortURL, err := sh.shortener.Shorten(shortenReq.URL)
	if err != nil {
		return liberr.InternalError(err.Error())
	}

	shortenResp := contract.ShortenResponse{ShortURL: shortURL}

	util.WriteSuccessResponse(http.StatusCreated, shortenResp, resp)
	return nil
}

func NewShortenHandler(shortener shortener.Shortener) *ShortenHandler {
	return &ShortenHandler{
		shortener: shortener,
	}
}
