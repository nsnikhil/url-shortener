package service

import (
	"go.uber.org/zap"
	"net/url"
	"urlshortner/pkg/shortener/consts"
	"urlshortner/pkg/shortener/contract"
	"urlshortner/pkg/shortener/store"
)

type ShortenerService interface {
	Shorten(req contract.ShortenRequest) (contract.ShortenResponse, error)
	Redirect(urlHash string) (string, error)
}

type defaultShortenerService struct {
	lgr       *zap.Logger
	str       *store.Store
	builder   URLBuilder
	generator HashGenerator
}

func (dss *defaultShortenerService) Shorten(req contract.ShortenRequest) (contract.ShortenResponse, error) {
	err := isValidURL(req.URL, dss.lgr)
	if err != nil {
		return contract.ShortenResponse{}, err
	}

	urlHash, err := dss.str.GetShortnerStore().GetURLHash(req.URL)
	if err == nil {
		return contract.ShortenResponse{ShortURL: dss.builder.Build(urlHash)}, nil
	}

	urlHash = dss.generator.Generate(req.URL)

	_, err = dss.str.GetShortnerStore().Save(req.URL, urlHash)
	if err != nil {
		return contract.ShortenResponse{}, err
	}

	return contract.ShortenResponse{ShortURL: dss.builder.Build(urlHash)}, nil
}

func (dss *defaultShortenerService) Redirect(urlHash string) (string, error) {
	longURL, err := dss.str.GetShortnerStore().GetURL(urlHash[1:])
	if err != nil {
		return consts.NA, err
	}

	return longURL, nil
}

func isValidURL(urlStr string, lgr *zap.Logger) error {
	_, err := url.Parse(urlStr)
	if err != nil {
		lgr.Error(err.Error())
		return err
	}

	return nil
}

func NewShortenerService(lgr *zap.Logger, generator HashGenerator, builder URLBuilder, str *store.Store) ShortenerService {
	return &defaultShortenerService{
		lgr:       lgr,
		generator: generator,
		builder:   builder,
		str:       str,
	}
}
