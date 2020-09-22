package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/http/internal/handler"
	"urlshortner/pkg/http/internal/middleware"
	"urlshortner/pkg/shortener"
)

func TestShortenHandler(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test shorten handler success",
			actualResult: func() (string, int) {
				req := contract.ShortenRequest{URL: "wikipedia.com"}
				b, err := json.Marshal(req)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(b))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}
				mockShortener.On("Shorten", "wikipedia.com").Return("sht.ly/abc", nil)

				sh := handler.NewShortenHandler(mockShortener)
				middleware.WithError(sh.Shorten)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "{\"data\":{\"short_url\":\"sht.ly/abc\"},\"error\":{},\"success\":true}",
			expectedCode:   http.StatusCreated,
		},
		{
			name: "test shorten handler fail when body is nil",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", nil)
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}

				sh := handler.NewShortenHandler(mockShortener)
				middleware.WithError(sh.Shorten)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "{\"error\":{\"code\":\"usx0001\",\"message\":\"request body is nil\"},\"success\":false}",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name: "test shorten handler fail when body is invalid",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString("invalid body"))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}

				sh := handler.NewShortenHandler(mockShortener)
				middleware.WithError(sh.Shorten)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "{\"error\":{\"code\":\"usx0001\",\"message\":\"invalid character 'i' looking for beginning of value\"},\"success\":false}",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name: "test shorten handler fail when service returns error",
			actualResult: func() (string, int) {
				req := contract.ShortenRequest{URL: "wikipedia.com"}
				b, err := json.Marshal(req)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(b))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}
				mockShortener.On("Shorten", "wikipedia.com").Return("", errors.New("fail to shorten url"))

				sh := handler.NewShortenHandler(mockShortener)
				middleware.WithError(sh.Shorten)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "{\"error\":{\"code\":\"usx0010\",\"message\":\"fail to shorten url\"},\"success\":false}",
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, cd := testCase.actualResult()

			assert.Equal(t, testCase.expectedCode, cd)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
