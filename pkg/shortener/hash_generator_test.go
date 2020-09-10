package shortener_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
	"testing"
	"urlshortner/pkg/shortener"
)

func TestHashGeneratorGenerate(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, error)
		expectedResult string
		expectedError  error
	}{
		{
			name: "test generate hash for a string",
			actualResult: func() (string, error) {
				hs := sha3.New512()
				length := 9

				gen := shortener.NewHashGenerator(hs, length)
				return gen.Generate("wikipedia.com")
			},
			expectedResult: "AreVAfnsk",
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

func TestHashGeneratorGenerateSameHash(t *testing.T) {
	testCases := []struct {
		name         string
		actualResult func() []string
	}{
		{
			"test generate same hash for a string",
			func() []string {
				hs := sha3.New512()
				length := 9
				gen := shortener.NewHashGenerator(hs, length)

				var res []string

				for i := 0; i < 10; i++ {
					hs, err := gen.Generate("wikipedia.com")
					require.NoError(t, err)

					res = append(res, hs)
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
				assert.Equal(t, curr, res[i])
				curr = res[i]
			}
		})
	}
}
