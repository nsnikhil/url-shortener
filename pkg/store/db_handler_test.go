package store_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/config"
	"urlshortner/pkg/store"
)

func TestGetDB(t *testing.T) {
	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test get db success",
			actualResult: func() error {
				handler := store.NewDBHandler(config.NewConfig().GetDatabaseConfig(), zap.NewNop())
				_, err := handler.GetDB()
				return err
			},
			expectedError: nil,
		},
		{
			name: "test get db failure invalid driver",
			actualResult: func() error {
				handler := store.NewDBHandler(config.DatabaseConfig{}, zap.NewNop())
				_, err := handler.GetDB()
				return err
			},
			expectedError: errors.New("sql: unknown driver \"\" (forgotten import?)"),
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
