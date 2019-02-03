package concourse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	sut := New("concourse", "http://concourse.example.com")
	assert.NotNil(t, sut)
}

func TestNew_PanicsIfBadURL(t *testing.T) {
	var testCases = []struct {
		URI     string
		Message string
	}{
		{URI: "hk>123", Message: "Malformed URL"},
		{URI: "", Message: "Empty URL"},
		{URI: "a", Message: "Relative URL is not valid"},
		{URI: "http://", Message: "No host present is not valid"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Message, func(st *testing.T) {
			assert.Panics(st, func() { New("concourse", testCase.URI) })
		})
	}
}
