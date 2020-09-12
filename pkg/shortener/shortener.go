package shortener

import (
	"go.uber.org/zap"
	"net/url"
	"urlshortner/pkg/store"
)

type Shortener interface {
	Shorten(url string) (string, error)
}

type defaultShortener struct {
	logger *zap.Logger

	store store.ShortenerStore

	builder   URLBuilder
	generator HashGenerator
}

func (ds *defaultShortener) Shorten(url string) (string, error) {
	err := isValidURL(url, ds.logger)
	if err != nil {
		return "", err
	}

	urlHash, err := ds.generator.Generate(url)
	if err != nil {
		return "", err
	}

	err = ds.store.Save(url, urlHash)
	if err != nil {
		return "", err
	}

	return ds.builder.Build(urlHash), nil
}

func isValidURL(urlStr string, lgr *zap.Logger) error {
	_, err := url.Parse(urlStr)
	if err != nil {
		lgr.Error(err.Error())
		return err
	}

	return nil
}

func NewShortener(lgr *zap.Logger, str store.ShortenerStore, builder URLBuilder, generator HashGenerator) Shortener {
	return &defaultShortener{
		logger:    lgr,
		store:     str,
		builder:   builder,
		generator: generator,
	}
}
