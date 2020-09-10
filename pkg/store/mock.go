package store

import "github.com/stretchr/testify/mock"

type MockShortnerStore struct {
	mock.Mock
}

func (mss *MockShortnerStore) Save(url, urlHash string) (int, error) {
	args := mss.Called(url, urlHash)
	return args.Int(0), args.Error(1)
}

func (mss *MockShortnerStore) GetURL(urlHash string) (string, error) {
	args := mss.Called(urlHash)
	return args.String(0), args.Error(1)
}

func (mss *MockShortnerStore) GetURLHash(url string) (string, error) {
	args := mss.Called(url)
	return args.String(0), args.Error(1)
}

func (mss *MockShortnerStore) Delete(url, urlHash string) (int, error) {
	args := mss.Called(url, urlHash)
	return args.Int(0), args.Error(1)
}
