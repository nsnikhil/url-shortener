package store_test

import (
	"database/sql"
	"errors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/config"
	"urlshortner/pkg/store"
)

func TestDatabaseSave(t *testing.T) {
	db, sdb := getDB(t)
	defer func() { _ = db.Close() }()
	defer cleanUp(t, db)

	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test save success",
			actualResult: func() error {
				return sdb.Save("wikipedia.com", "some-random-hash")
			},
		},
		{
			name: "test save failure when urlhash is empty",
			actualResult: func() error {
				return sdb.Save("some-url.com", "")
			},
			expectedError: errors.New("pq: new row for relation \"shortener\" violates check constraint \"shortener_urlhash_check\""),
		},
		{
			name: "test save failure when url is empty",
			actualResult: func() error {
				return sdb.Save("", "some-random-hash")
			},
			expectedError: errors.New("pq: new row for relation \"shortener\" violates check constraint \"shortener_url_check\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), testCase.actualResult().Error())
			} else {
				assert.Nil(t, testCase.actualResult())
			}
		})
	}
}

func TestShortenerDatabaseGet(t *testing.T) {
	db, sdb := getDB(t)
	defer func() { _ = db.Close() }()
	defer cleanUp(t, db)

	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test shortener get success",
			actualResult: func() (string, error) {
				err := sdb.Save("wikipedia.com", "some-random-hash")
				require.NoError(t, err)
				return sdb.Get("some-random-hash")
			},
			expectedResult: "wikipedia.com",
		},
		{
			name: "test shortener get failure",
			actualResult: func() (string, error) {
				return sdb.Get("not-present-hash")
			},
			expectedResult: "",
			expectedError:  errors.New("sql: no rows in result set"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}

func getDB(t *testing.T) (*sql.DB, store.ShortenerDatabase) {
	cfg := config.NewConfig()
	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), zap.NewNop())

	db, err := dbHandler.GetDB()
	require.NoError(t, err)

	nr, _ := newrelic.NewApplication()
	return db, store.NewShortenerDatabase(db, zap.NewNop(), nr)
}

func cleanUp(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE shortener")
	require.NoError(t, err)
}
