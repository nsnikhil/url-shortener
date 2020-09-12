package store

import "github.com/stretchr/testify/mock"

type MockShortenerDatabase struct {
	mock.Mock
}

func (msd *MockShortenerDatabase) Save(url, urlHash string) error {
	args := msd.Called(url, urlHash)
	return args.Error(0)
}

func (msd *MockShortenerDatabase) Get(urlHash string) (string, error) {
	args := msd.Called(urlHash)
	return args.String(0), args.Error(1)
}

type MockShortenerCache struct {
	mock.Mock
}

func (msd *MockShortenerCache) Save(url, urlHash string, ttl int) error {
	args := msd.Called(url, urlHash, ttl)
	return args.Error(0)
}

func (msd *MockShortenerCache) Get(urlHash string) (string, error) {
	args := msd.Called(urlHash)
	return args.String(0), args.Error(1)
}

type MockShortenerStore struct {
	mock.Mock
}

func (mss *MockShortenerStore) Save(url, urlHash string) error {
	args := mss.Called(url, urlHash)
	return args.Error(1)
}

func (mss *MockShortenerStore) GetURL(urlHash string) (string, error) {
	args := mss.Called(urlHash)
	return args.String(0), args.Error(1)
}
