package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"urlshortner/pkg/shortener/contract"
	"urlshortner/pkg/shortener/service"
	"urlshortner/pkg/shortener/store"
)

func TestShortenerServiceShorten(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (contract.ShortenResponse, error)
		expectedResult contract.ShortenResponse
		expectedError  error
	}{
		{
			name: "test shorten new long url",
			actualResult: func() (contract.ShortenResponse, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockShortenerStore.On("GetURLHash", "veryLongURL.com").Return("", errors.New("not found"))
				mockShortenerStore.On("Save", "veryLongURL.com", "MLReNfDWL").Return(1, nil)

				mockURLBuilder := &service.MockURLBuilder{}
				mockURLBuilder.On("Build", "MLReNfDWL").Return("sht.ly/MLReNfDWL")

				mockHashGenerator := &service.MockHashGenerator{}
				mockHashGenerator.On("Generate", "veryLongURL.com").Return("MLReNfDWL")

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLBuilder, store.NewStore(mockShortenerStore))

				req := contract.ShortenRequest{URL: "veryLongURL.com"}

				return ss.Shorten(req)
			},
			expectedResult: contract.ShortenResponse{
				ShortURL: "sht.ly/MLReNfDWL",
			},
		},
		{
			name: "test shorten existing long url",
			actualResult: func() (contract.ShortenResponse, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockShortenerStore.On("GetURLHash", "veryLongURL.com").Return("MLReNfDWL", nil)

				mockURLBuilder := &service.MockURLBuilder{}
				mockURLBuilder.On("Build", "MLReNfDWL").Return("sht.ly/MLReNfDWL")

				mockHashGenerator := &service.MockHashGenerator{}

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLBuilder, store.NewStore(mockShortenerStore))

				req := contract.ShortenRequest{URL: "veryLongURL.com"}

				return ss.Shorten(req)
			},
			expectedResult: contract.ShortenResponse{
				ShortURL: "sht.ly/MLReNfDWL",
			},
		},
		{
			name: "test shorten failure when url is invalid",
			actualResult: func() (contract.ShortenResponse, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockURLBuilder := &service.MockURLBuilder{}
				mockHashGenerator := &service.MockHashGenerator{}

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLBuilder, store.NewStore(mockShortenerStore))

				req := contract.ShortenRequest{URL: "#@$%^"}

				return ss.Shorten(req)
			},
			expectedResult: contract.ShortenResponse{},
			expectedError:  errors.New("parse \"#@$%^\": invalid URL escape \"%^\""),
		},
		{
			name: "test shorten failure when save fails",
			actualResult: func() (contract.ShortenResponse, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockShortenerStore.On("GetURLHash", "veryLongURL.com").Return("", errors.New("not found"))
				mockShortenerStore.On("Save", "veryLongURL.com", "MLReNfDWL").Return(0, errors.New("failed to save"))

				mockURLBuilder := &service.MockURLBuilder{}
				mockURLBuilder.On("Build", "MLReNfDWL").Return("sht.ly/MLReNfDWL")

				mockHashGenerator := &service.MockHashGenerator{}
				mockHashGenerator.On("Generate", "veryLongURL.com").Return("MLReNfDWL")

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLBuilder, store.NewStore(mockShortenerStore))

				req := contract.ShortenRequest{URL: "veryLongURL.com"}

				return ss.Shorten(req)
			},
			expectedResult: contract.ShortenResponse{},
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

func TestShortenerServiceRedirect(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test redirect success",
			actualResult: func() (string, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockShortenerStore.On("GetURL", "MLReNfDWL").Return("veryLongURL.com", nil)

				mockURLGenerator := &service.MockURLBuilder{}
				mockHashGenerator := &service.MockHashGenerator{}

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLGenerator, store.NewStore(mockShortenerStore))
				return ss.Redirect("/MLReNfDWL")
			},
			expectedResult: "veryLongURL.com",
		},
		{
			name: "test redirect failure",
			actualResult: func() (string, error) {
				mockShortenerStore := &store.MockShortnerStore{}
				mockShortenerStore.On("GetURL", "MLReNfDWL").Return("na", errors.New("not found"))

				mockURLGenerator := &service.MockURLBuilder{}
				mockHashGenerator := &service.MockHashGenerator{}

				ss := service.NewShortenerService(zap.NewNop(), mockHashGenerator, mockURLGenerator, store.NewStore(mockShortenerStore))
				return ss.Redirect("/MLReNfDWL")
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
