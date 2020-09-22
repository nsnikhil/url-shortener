package store_test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/config"
	"urlshortner/pkg/store"
)

func TestGetCache(t *testing.T) {
	testCases := []struct {
		name         string
		actualResult func() error
		hasError     bool
	}{
		{
			name: "test get cache success",
			actualResult: func() error {
				handler := store.NewCacheHandler(config.NewConfig().GetRedisConfig(), zap.NewNop())
				_, err := handler.GetCache()
				return err
			},
			hasError: false,
		},
		{
			name: "test get cache failure",
			actualResult: func() error {
				handler := store.NewCacheHandler(config.RedisConfig{}, zap.NewNop())
				_, err := handler.GetCache()
				return err
			},
			hasError:      true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.hasError {
				assert.NotNil(t, testCase.actualResult())
			} else {
				assert.Nil(t, testCase.actualResult())
			}
		})
	}
}
