package handler

import (
	"net/http"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/liberr"
)

const locationHeader = "Location"

type RedirectHandler struct {
	elongator elongator.Elongator
}

func (rh *RedirectHandler) Redirect(resp http.ResponseWriter, req *http.Request) error {
	longURL, err := rh.elongator.Elongate(req.URL.Path[1:])
	if err != nil {
		return liberr.InternalError(err.Error())
	}

	resp.Header().Set(locationHeader, longURL)
	resp.WriteHeader(http.StatusMovedPermanently)

	return nil
}

func NewRedirectHandler(elongator elongator.Elongator) *RedirectHandler {
	return &RedirectHandler{
		elongator: elongator,
	}
}
