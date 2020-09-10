package handler_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/handler"
)

func TestRedirectHandler(t *testing.T) {
	testCases := []struct {
		name               string
		actualResult       func() (http.Header, int)
		expectedHeader     http.Header
		expectedStatusCode int
	}{
		{
			name: "test redirect handler success",
			actualResult: func() (http.Header, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/abc", nil)
				require.NoError(t, err)

				mockElongator := &elongator.MockElongator{}
				mockElongator.On("Elongate", "abc").Return("veryLongUrl.com", nil)

				handler.RedirectHandler(zap.NewNop(), &statsd.Client{}, mockElongator)(w, r)

				return w.Header(), w.Code
			},
			expectedHeader:     http.Header{"Location": []string{"veryLongUrl.com"}},
			expectedStatusCode: http.StatusMovedPermanently,
		},
		{
			name: "test redirect handler failure",
			actualResult: func() (http.Header, int) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodPost, "/abc", nil)
				require.NoError(t, err)

				mockElongator := &elongator.MockElongator{}
				mockElongator.On("Elongate", "abc").Return("na", errors.New("invalid hash"))

				handler.RedirectHandler(zap.NewNop(), &statsd.Client{}, mockElongator)(w, r)

				return w.Header(), w.Code
			},
			expectedHeader:     http.Header{},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hd, cd := testCase.actualResult()

			assert.Equal(t, testCase.expectedStatusCode, cd)
			assert.Equal(t, testCase.expectedHeader, hd)
		})
	}
}
