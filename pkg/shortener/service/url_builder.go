package service

import (
	"fmt"
)

type URLBuilder interface {
	Build(hash string) string
}

type defaultURLBuilder struct {
	baseURL string
}

func (dug *defaultURLBuilder) Build(hash string) string {
	return fmt.Sprintf("%s/%s", dug.baseURL, hash)
}

func NewURLBuilder(baseURL string) URLBuilder {
	return &defaultURLBuilder{
		baseURL: baseURL,
	}
}
