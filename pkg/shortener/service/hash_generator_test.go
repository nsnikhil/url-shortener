package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHashGeneratorGenerateDifferentHash(t *testing.T) {
	testCases := []struct {
		name         string
		actualResult func() []string
	}{
		{
			name: "test generate different hash for a string",
			actualResult: func() []string {
				hs := sha3.New512()
				gen := service.NewHashGenerator(hs, 9)

				var res []string

				for i := 0; i < 10; i++ {
					res = append(res, gen.Generate("wikipedia.com"))
				}

				return res
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := testCase.actualResult()
			sz := len(res)
			require.True(t, sz >= 2)

			curr := res[0]
			for i := 1; i < sz; i++ {
				assert.False(t, curr == res[i])
				curr = res[i]
			}
		})
	}
}
