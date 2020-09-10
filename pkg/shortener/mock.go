package shortener

import "github.com/stretchr/testify/mock"

type MockHashGenerator struct {
	mock.Mock
}

func (mhg *MockHashGenerator) Generate(str string) (string, error) {
	args := mhg.Called(str)
	return args.String(0), args.Error(1)
}

type MockURLBuilder struct {
	mock.Mock
}

func (mub *MockURLBuilder) Build(hash string) string {
	args := mub.Called(hash)
	return args.String(0)
}

type MockShortener struct {
	mock.Mock
}

func (ms *MockShortener) Shorten(url string) (string, error) {
	args := ms.Called(url)
	return args.String(0), args.Error(1)
}
