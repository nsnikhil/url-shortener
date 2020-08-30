package store_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
	"urlshortner/cmd/config"
	"urlshortner/pkg/shortener/store"
)

func TestShortnerStoreSave(t *testing.T) {
	cfg := config.NewConfig()
	lgr := zap.NewNop()

	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), lgr)
	cacheHandler := store.NewCacheHandler(cfg.GetRedisConfig(), lgr)

	db, err := dbHandler.GetDB()
	require.NoError(t, err)

	cache, err := cacheHandler.GetCache()
	require.NoError(t, err)

	testCases := []struct {
		name           string
		actualResult   func() (int, error)
		expectedResult int
		expectedError  error
	}{
		{
			name: "test save single record in db and cache",
			actualResult: func() (int, error) {
				ss := store.NewShortnerStore(cache, db, lgr)
				c, err := ss.Save("verLongURL.com", "sht.ly/abc")
				truncateStore(t, db, cache)
				return c, err
			},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name: "test save single record in db and cache",
			actualResult: func() (int, error) {
				ss := store.NewShortnerStore(cache, db, lgr)
				res := 0

				c, err := ss.Save("verLongURL.com", "sht.ly/abc")
				require.NoError(t, err)
				res += c

				c, err = ss.Save("otherURL.com", "sht.ly/def")
				require.NoError(t, err)
				res += c

				c, err = ss.Save("thisURL.com", "sht.ly/ghi")
				require.NoError(t, err)
				res += c

				c, err = ss.Save("thatURL.com", "sht.ly/jkl")
				require.NoError(t, err)
				res += c

				truncateStore(t, db, cache)
				return res, err
			},
			expectedResult: 4,
			expectedError:  nil,
		},
		{
			name: "test save single failure when url is empty",
			actualResult: func() (int, error) {
				ss := store.NewShortnerStore(cache, db, lgr)
				c, err := ss.Save("", "sht.ly/abc")
				truncateStore(t, db, cache)
				return c, err
			},
			expectedResult: 0,
			expectedError:  errors.New("pq: new row for relation \"shortener\" violates check constraint \"shortener_url_check\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedResult, res)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestShortnerStoreGetURL(t *testing.T) {
	cfg := config.NewConfig()
	lgr := zap.NewNop()

	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), lgr)
	cacheHandler := store.NewCacheHandler(cfg.GetRedisConfig(), lgr)

	db, err := dbHandler.GetDB()
	require.NoError(t, err)

	cache, err := cacheHandler.GetCache()
	require.NoError(t, err)

	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test fetch record from cache",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				_, err := ss.Save("verLongURL.com", "randomHash")
				require.NoError(t, err)
				time.Sleep(time.Millisecond)

				longURL, err := ss.GetURL("randomHash")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "verLongURL.com",
		},
		{
			name: "test fetch record from db",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				_, err := ss.Save("verLongURL.com", "randomHash")
				require.NoError(t, err)

				time.Sleep(time.Millisecond)
				require.NoError(t, cache.Del(context.Background(), "randomHash").Err())

				longURL, err := ss.GetURL("randomHash")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "verLongURL.com",
		},
		{
			name: "test fetch record failure",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				longURL, err := ss.GetURL("randomHash")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "na",
			expectedError:  errors.New("sql: no rows in result set"),
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

func TestShortnerStoreGetURLHash(t *testing.T) {
	cfg := config.NewConfig()
	lgr := zap.NewNop()

	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), lgr)
	cacheHandler := store.NewCacheHandler(cfg.GetRedisConfig(), lgr)

	db, err := dbHandler.GetDB()
	require.NoError(t, err)

	cache, err := cacheHandler.GetCache()
	require.NoError(t, err)

	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test fetch record from cache",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				_, err := ss.Save("verLongURL.com", "randomHash")
				require.NoError(t, err)
				time.Sleep(time.Millisecond)

				longURL, err := ss.GetURLHash("verLongURL.com")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "randomHash",
		},
		{
			name: "test fetch record from db",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				_, err := ss.Save("verLongURL.com", "randomHash")
				require.NoError(t, err)

				time.Sleep(time.Millisecond)
				require.NoError(t, cache.Del(context.Background(), "verLongURL.com").Err())

				longURL, err := ss.GetURLHash("verLongURL.com")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "randomHash",
		},
		{
			name: "test fetch record failure",
			actualResult: func() (string, error) {
				ss := store.NewShortnerStore(cache, db, lgr)

				longURL, err := ss.GetURLHash("verLongURL.com")

				truncateStore(t, db, cache)

				return longURL, err
			},
			expectedResult: "na",
			expectedError:  errors.New("sql: no rows in result set"),
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

func truncateStore(t *testing.T, db *sql.DB, cache *redis.Client) {
	_, err := db.Exec("TRUNCATE shortener")
	require.NoError(t, err)
	require.NoError(t, cache.FlushAll(context.Background()).Err())
}
