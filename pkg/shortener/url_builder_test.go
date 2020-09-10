package shortener_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"urlshortner/pkg/shortener"
)

func TestURLGeneratorGenerate(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() string
		expectedResult string
	}{
		{
			name: "test generate short url",
			actualResult: func() string {
				baseURL := "sht.ly"

				builder := shortener.NewURLBuilder(baseURL)

				return builder.Build("MLReNfDWL")
			},
			expectedResult: "sht.ly/MLReNfDWL",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}
