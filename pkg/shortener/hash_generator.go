package shortener

import (
	"encoding/base64"
	"hash"
	"strings"
)

type HashGenerator interface {
	Generate(url string) (string, error)
}

type urlHashGenerator struct {
	hash   hash.Hash
	length int
}

func (uhg *urlHashGenerator) Generate(url string) (string, error) {
	hashCode, err := generateHash(uhg.hash, url)
	if err != nil {
		return "", err
	}

	return format(encode(hashCode), uhg.length), nil
}

func generateHash(hs hash.Hash, url string) ([]byte, error) {
	hs.Reset()

	_, err := hs.Write([]byte(url))
	if err != nil {
		return nil, err
	}

	return hs.Sum(nil), nil
}

func encode(bytes []byte) string {
	return base64.URLEncoding.EncodeToString(bytes)
}

func format(hash string, length int) string {
	bl := strings.Builder{}
	sz := len(hash)

	for i := 0; i < sz; i++ {
		if bl.Len() >= length {
			break
		}

		if hash[i] >= 65 && hash[i] <= 90 || hash[i] >= 97 && hash[i] <= 122 {
			bl.WriteRune(rune(hash[i]))
		}
	}

	return bl.String()
}

func NewHashGenerator(hash hash.Hash, length int) HashGenerator {
	return &urlHashGenerator{
		hash:   hash,
		length: length,
	}
}
