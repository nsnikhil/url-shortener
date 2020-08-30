package service_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"urlshortner/pkg/shortener/service"
)

func TestGetShortnerSService(t *testing.T) {
	ss := &service.MockShortenerService{}
	st := service.NewService(ss)
	assert.Equal(t, ss, st.GetShortenerService())
}
