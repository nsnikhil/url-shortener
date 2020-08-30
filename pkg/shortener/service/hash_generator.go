package service

import (
	"encoding/base64"
	"hash"
	"strings"
)

type HashGenerator interface {
	Generate(str string) string
}

type defaultHashGenerator struct {
	hash   hash.Hash
	length int
}

func (dhg *defaultHashGenerator) Generate(str string) string {
	_, _ = dhg.hash.Write([]byte(str))
	sha := base64.URLEncoding.EncodeToString(dhg.hash.Sum(nil))
	return format(sha, dhg.length)
}

func format(sha string, length int) string {
	bl := strings.Builder{}
	sz := len(sha)

	for i := 0; i < sz; i++ {
		if bl.Len() >= length {
			break
		}

		if sha[i] >= 65 && sha[i] <= 90 || sha[i] >= 97 && sha[i] <= 122 {
			bl.WriteRune(rune(sha[i]))
		}
	}

	return bl.String()
}

func NewHashGenerator(hash hash.Hash, length int) HashGenerator {
	return &defaultHashGenerator{
		hash:   hash,
		length: length,
	}
}
