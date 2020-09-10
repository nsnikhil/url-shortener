package shortener

import (
	"fmt"
)

type URLBuilder interface {
	Build(hash string) string
}

type shortURLBuilder struct {
	baseURL string
}

func (sur *shortURLBuilder) Build(hash string) string {
	return fmt.Sprintf("%s/%s", sur.baseURL, hash)
}

func NewURLBuilder(baseURL string) URLBuilder {
	return &shortURLBuilder{
		baseURL: baseURL,
	}
}
