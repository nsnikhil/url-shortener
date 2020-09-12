package store_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/store"
)

func TestShortnerStoreSave(t *testing.T) {
	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test save record",
			actualResult: func() error {
				url := "wikipedia.com"
				urlHash := "some-random-hash"

				db := &store.MockShortenerDatabase{}
				db.On("Save", url, urlHash).Return(nil)

				cache := &store.MockShortenerCache{}
				cache.On("Save", url, urlHash, 3600).Return(nil)

				ss := store.NewShortenerStore(db, cache, zap.NewNop())

				return ss.Save(url, urlHash)
			},
			expectedError: nil,
		},
		{
			name: "test save failure db return error",
			actualResult: func() error {
				url := "some-url"
				urlHash := "some-random-hash"

				db := &store.MockShortenerDatabase{}
				db.On("Save", url, urlHash).Return(errors.New("url is empty"))

				cache := &store.MockShortenerCache{}

				ss := store.NewShortenerStore(db, cache, zap.NewNop())

				return ss.Save(url, urlHash)
			},
			expectedError: errors.New("url is empty"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedError, testCase.actualResult())
		})
	}
}

func TestShortnerStoreGetURL(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test get record from cache",
			actualResult: func() (string, error) {
				url := "wikipedia.com"
				urlHash := "some-random-hash"

				db := &store.MockShortenerDatabase{}

				cache := &store.MockShortenerCache{}
				cache.On("Get", urlHash).Return(url, nil)

				ss := store.NewShortenerStore(db, cache, zap.NewNop())

				return ss.GetURL(urlHash)
			},
			expectedResult: "wikipedia.com",
			expectedError:  nil,
		},
		{
			name: "test get record from db",
			actualResult: func() (string, error) {
				url := "wikipedia.com"
				urlHash := "some-random-hash"

				db := &store.MockShortenerDatabase{}
				db.On("Get", urlHash).Return(url, nil)

				cache := &store.MockShortenerCache{}
				cache.On("Get", urlHash).Return("", errors.New("some error"))

				ss := store.NewShortenerStore(db, cache, zap.NewNop())

				return ss.GetURL(urlHash)
			},
			expectedResult: "wikipedia.com",
			expectedError:  nil,
		},
		{
			name: "test get record failure",
			actualResult: func() (string, error) {
				urlHash := "some-random-hash"

				db := &store.MockShortenerDatabase{}
				db.On("Get", urlHash).Return("", errors.New("some error"))

				cache := &store.MockShortenerCache{}
				cache.On("Get", urlHash).Return("", errors.New("some error"))

				ss := store.NewShortenerStore(db, cache, zap.NewNop())

				return ss.GetURL(urlHash)
			},
			expectedResult: "",
			expectedError:  errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedResult, res)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
