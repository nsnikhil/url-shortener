package elongator_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/store"
)

func TestElongatorElongate(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test redirect success",
			actualResult: func() (string, error) {
				mockStore := &store.MockShortenerStore{}
				mockStore.On("GetURL", "MLReNfDWL").Return("wikipedia.com", nil)

				el := elongator.NewElongator(zap.NewNop(), mockStore)

				return el.Elongate("MLReNfDWL")
			},
			expectedResult: "wikipedia.com",
		},
		{
			name: "test redirect failure",
			actualResult: func() (string, error) {
				mockStore := &store.MockShortenerStore{}
				mockStore.On("GetURL", "MLReNfDWL").Return("na", errors.New("not found"))

				el := elongator.NewElongator(zap.NewNop(), mockStore)

				return el.Elongate("MLReNfDWL")
			},
			expectedResult: "na",
			expectedError:  errors.New("not found"),
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
