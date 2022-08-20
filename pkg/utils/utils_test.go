package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseInt(t *testing.T) {
	assert.Equal(t, int64(1234567890), ParseInt("1.234.567.890"))
}

func TestToInt(t *testing.T) {
	assert.Equal(t, 1234567890, ToInt([]byte("1234567890")))
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, int64(2), MinInt(5, 2, 3))
}
