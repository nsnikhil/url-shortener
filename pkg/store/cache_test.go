package store_test

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
	"urlshortner/pkg/config"
	"urlshortner/pkg/store"
)

func TestShortenerCacheSave(t *testing.T) {
	url := "wikipedia.com"
	urlHash := "some-hash-value"
	ttlInSec := 1

	cl, cache := getRedis(t)
	defer func() { _ = cl.Close() }()

	err := cache.Save(url, urlHash, ttlInSec)
	require.NoError(t, err)
}

func TestUrlShortenerCacheGet(t *testing.T) {
	cl, cache := getRedis(t)
	defer func() { _ = cl.Close() }()

	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test get from cache success",
			actualResult: func() (string, error) {
				urlHash := "some-hash-value"
				err := cache.Save("wikipedia.com", urlHash, 1)
				require.NoError(t, err)

				return cache.Get(urlHash)
			},
			expectedResult: "wikipedia.com",
		},
		{
			name: "test get fails when key is not present",
			actualResult: func() (string, error) {
				return cache.Get("invalid-key")
			},
			expectedResult: "",
			expectedError:  errors.New("redis: nil"),
		},
		{
			name: "test get fails when ttl expires",
			actualResult: func() (string, error) {
				urlHash := "some-hash-value"
				err := cache.Save("wikipedia.com", urlHash, 1)
				require.NoError(t, err)

				time.Sleep(time.Second)
				return cache.Get(urlHash)
			},
			expectedResult: "",
			expectedError:  errors.New("redis: nil"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}

			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}

func getRedis(t *testing.T) (*redis.Client, store.ShortenerCache) {
	cfg := config.NewConfig()
	lgr := zap.NewNop()

	ch := store.NewCacheHandler(cfg.GetRedisConfig(), lgr)

	cl, err := ch.GetCache()
	require.NoError(t, err)

	return cl, store.NewShortenerCache(cl, lgr)
}
