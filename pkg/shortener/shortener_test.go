package shortener_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/shortener"
	"urlshortner/pkg/store"
)

func TestShortenerServiceShorten(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test shorten new long url",
			actualResult: func() (string, error) {
				urlHash := "MLReNfDWL"
				shortURL := "sht.ly/MLReNfDWL"
				longURL := "wikipedia.com"

				mockStore := &store.MockShortnerStore{}
				mockStore.On("Save", longURL, urlHash).Return(1, nil)

				mockHashGenerator := &shortener.MockHashGenerator{}
				mockHashGenerator.On("Generate", longURL).Return(urlHash, nil)

				mockURLBuilder := &shortener.MockURLBuilder{}
				mockURLBuilder.On("Build", urlHash).Return(shortURL)

				sht := shortener.NewShortener(zap.NewNop(), store.NewStore(mockStore), mockURLBuilder, mockHashGenerator)

				return sht.Shorten(longURL)
			},
			expectedResult: "sht.ly/MLReNfDWL",
		},
		{
			name: "test shorten existing long url",
			actualResult: func() (string, error) {
				urlHash := "MLReNfDWL"
				shortURL := "sht.ly/MLReNfDWL"
				longURL := "wikipedia.com"

				mockStore := &store.MockShortnerStore{}
				mockStore.On("Save", longURL, urlHash).Return(1, nil)

				mockHashGenerator := &shortener.MockHashGenerator{}
				mockHashGenerator.On("Generate", longURL).Return(urlHash, nil)

				mockURLBuilder := &shortener.MockURLBuilder{}
				mockURLBuilder.On("Build", urlHash).Return(shortURL)

				sht := shortener.NewShortener(zap.NewNop(), store.NewStore(mockStore), mockURLBuilder, mockHashGenerator)

				return sht.Shorten(longURL)
			},
			expectedResult: "sht.ly/MLReNfDWL",
		},
		{
			name: "test shorten failure when url is invalid",
			actualResult: func() (string, error) {
				mockStore := &store.MockShortnerStore{}
				mockURLBuilder := &shortener.MockURLBuilder{}
				mockHashGenerator := &shortener.MockHashGenerator{}

				sht := shortener.NewShortener(zap.NewNop(), store.NewStore(mockStore), mockURLBuilder, mockHashGenerator)

				return sht.Shorten("#@$%^")
			},
			expectedResult: "",
			expectedError:  errors.New("parse \"#@$%^\": invalid URL escape \"%^\""),
		},
		{
			name: "test shorten failure when generate return error",
			actualResult: func() (string, error) {
				longURL := "wikipedia.com"

				mockStore := &store.MockShortnerStore{}
				mockURLBuilder := &shortener.MockURLBuilder{}
				mockHashGenerator := &shortener.MockHashGenerator{}
				mockHashGenerator.On("Generate", longURL).Return("", errors.New("failed to generate hash"))

				sht := shortener.NewShortener(zap.NewNop(), store.NewStore(mockStore), mockURLBuilder, mockHashGenerator)

				return sht.Shorten("wikipedia.com")
			},
			expectedResult: "",
			expectedError:  errors.New("failed to generate hash"),
		},
		{
			name: "test shorten failure when save fails",
			actualResult: func() (string, error) {
				urlHash := "MLReNfDWL"
				shortURL := "sht.ly/MLReNfDWL"
				longURL := "wikipedia.com"

				mockStore := &store.MockShortnerStore{}
				mockStore.On("Save", longURL, urlHash).Return(0, errors.New("failed to save"))

				mockHashGenerator := &shortener.MockHashGenerator{}
				mockHashGenerator.On("Generate", longURL).Return(urlHash, nil)

				mockURLBuilder := &shortener.MockURLBuilder{}
				mockURLBuilder.On("Build", urlHash).Return(shortURL)

				sht := shortener.NewShortener(zap.NewNop(), store.NewStore(mockStore), mockURLBuilder, mockHashGenerator)

				return sht.Shorten(longURL)
			},
			expectedResult: "",
			expectedError:  errors.New("failed to save"),
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
