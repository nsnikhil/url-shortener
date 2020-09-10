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
	str *store.Store
}

func (de *defaultElongator) Elongate(hash string) (string, error) {
	longURL, err := de.str.GetShortnerStore().GetURL(hash)
	if err != nil {
		return na, err
	}

	return longURL, nil
}

func NewElongator(lgr *zap.Logger, str *store.Store) Elongator {
	return &defaultElongator{
		lgr: lgr,
		str: str,
	}
}
