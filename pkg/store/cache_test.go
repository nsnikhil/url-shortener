package store_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/config"
	"urlshortner/pkg/store"
)

func TestGetCache(t *testing.T) {
	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test get cache success",
			actualResult: func() error {
				handler := store.NewCacheHandler(config.NewConfig().GetRedisConfig(), zap.NewNop())
				_, err := handler.GetCache()
				return err
			},
		},
		{
			name: "test get cache failure",
			actualResult: func() error {
				handler := store.NewCacheHandler(config.RedisConfig{}, zap.NewNop())
				_, err := handler.GetCache()
				return err
			},
			expectedError: errors.New("dial tcp :0: connect: can't assign requested address"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expectedError == nil {
				assert.Equal(t, testCase.expectedError, testCase.actualResult())
			} else {
				assert.Equal(t, testCase.expectedError.Error(), testCase.actualResult().Error())
			}
		})
	}
}
