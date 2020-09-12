package elongator

import (
	"go.uber.org/zap"
	"urlshortner/pkg/store"
)

const na = "na"

type Elongator interface {
	Elongate(hash string) (string, error)
}

type defaultElongator struct {
	lgr *zap.Logger
	str store.ShortenerStore
}

func (de *defaultElongator) Elongate(hash string) (string, error) {
	longURL, err := de.str.GetURL(hash)
	if err != nil {
		return na, err
	}

	return longURL, nil
}

func NewElongator(lgr *zap.Logger, str store.ShortenerStore) Elongator {
	return &defaultElongator{
		lgr: lgr,
		str: str,
	}
}
