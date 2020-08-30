package service

import (
	"github.com/stretchr/testify/mock"
	"urlshortner/pkg/shortener/contract"
)

type MockHashGenerator struct {
	mock.Mock
}

func (mhg *MockHashGenerator) Generate(str string) string {
	args := mhg.Called(str)
	return args.String(0)
}

type MockURLBuilder struct {
	mock.Mock
}

func (mug *MockURLBuilder) Build(hash string) string {
	args := mug.Called(hash)
	return args.String(0)
}

type MockShortenerService struct {
	mock.Mock
}

func (mss *MockShortenerService) Shorten(req contract.ShortenRequest) (contract.ShortenResponse, error) {
	args := mss.Called(req)
	return args.Get(0).(contract.ShortenResponse), args.Error(1)
}

func (mss *MockShortenerService) Redirect(urlHash string) (string, error) {
	args := mss.Called(urlHash)
	return args.String(0), args.Error(1)
}
