package liberr_test

import (
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
	"urlshortner/pkg/http/liberr"
)

func TestGenericErrorGetErrorCode(t *testing.T) {
	ge := liberr.NewResponseError("us-def", http.StatusBadRequest, "some reason")

	assert.Equal(t, "us-def", ge.ErrorCode())
	assert.Equal(t, http.StatusBadRequest, ge.StatusCode())
	assert.Equal(t, "some reason", ge.Error())
}
