package store_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"urlshortner/pkg/shortener/store"
)

func TestGetShortnerStore(t *testing.T) {
	ss := &store.MockShortnerStore{}
	st := store.NewStore(ss)
	assert.Equal(t, ss, st.GetShortnerStore())
}
