package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/http/handler"
	"urlshortner/pkg/reporters"
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
				req := contract.ShortenRequest{URL: "veryLongUrl.com"}
				b, err := json.Marshal(req)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(b))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}
				mockShortener.On("Shorten", "veryLongUrl.com").Return("sht.ly/abc", nil)

				mockStatsD := &reporters.MockStatsDClient{}
				mockStatsD.On("ReportAttempt", "shorten")
				mockStatsD.On("ReportSuccess", "shorten")

				handler.ShortenHandler(zap.NewNop(), mockStatsD, mockShortener)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "{\"short_url\":\"sht.ly/abc\"}",
			expectedCode:   http.StatusOK,
		},
		{
			name: "test shorten handler fail when body is nil",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", nil)
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}

				mockStatsD := &reporters.MockStatsDClient{}
				mockStatsD.On("ReportAttempt", "shorten")
				mockStatsD.On("ReportFailure", "shorten")

				handler.ShortenHandler(zap.NewNop(), mockStatsD, mockShortener)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "body is nil",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name: "test shorten handler fail when body is invalid",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString("invalid body"))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}

				mockStatsD := &reporters.MockStatsDClient{}
				mockStatsD.On("ReportAttempt", "shorten")
				mockStatsD.On("ReportFailure", "shorten")

				handler.ShortenHandler(zap.NewNop(), mockStatsD, mockShortener)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "invalid character 'i' looking for beginning of value",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name: "test shorten handler fail when service returns error",
			actualResult: func() (string, int) {
				req := contract.ShortenRequest{URL: "veryLongUrl.com"}
				b, err := json.Marshal(req)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(b))
				require.NoError(t, err)

				mockShortener := &shortener.MockShortener{}
				mockShortener.On("Shorten", "veryLongUrl.com").Return("", errors.New("fail to shorten url"))

				mockStatsD := &reporters.MockStatsDClient{}
				mockStatsD.On("ReportAttempt", "shorten")
				mockStatsD.On("ReportFailure", "shorten")

				handler.ShortenHandler(zap.NewNop(), mockStatsD, mockShortener)(w, r)

				return w.Body.String(), w.Code
			},
			expectedResult: "fail to shorten url",
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
