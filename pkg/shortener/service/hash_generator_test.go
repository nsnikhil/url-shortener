package service_test

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"testing"
	"urlshortner/pkg/shortener/service"
)

func TestHashGeneratorGenerate(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() string
		expectedResult string
	}{
		{
			name: "test generate hash for a string",
			actualResult: func() string {
				hs := sha3.New512()
				gen := service.NewHashGenerator(hs, 9)
				return gen.Generate("veryLongURL.com")
			},
			expectedResult: "MLReNfDWL",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}
