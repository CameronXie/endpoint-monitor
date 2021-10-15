package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestError_Error(t *testing.T) {
	a := assert.New(t)
	a.Equal(
		"something went wrong",
		(&RequestError{URL: "https://example.com", Message: "something went wrong"}).Error(),
	)
}
