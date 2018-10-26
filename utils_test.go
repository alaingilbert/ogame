package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	assert.Equal(t, 1234567890, ParseInt("1.234.567.890"))
}

func TestToInt(t *testing.T) {
	assert.Equal(t, 1234567890, toInt([]byte("1234567890")))
}
