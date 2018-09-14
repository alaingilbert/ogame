package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_String(t *testing.T) {
	assert.Equal(t, "[1:2:3]", Coordinate{1, 2, 3}.String())
}

func TestCoordinate_Equal(t *testing.T) {
	assert.True(t, Coordinate{1, 2, 3}.Equal(Coordinate{1, 2, 3}))
	assert.False(t, Coordinate{1, 2, 3}.Equal(Coordinate{2, 2, 3}))
	assert.False(t, Coordinate{1, 2, 3}.Equal(Coordinate{1, 3, 3}))
	assert.False(t, Coordinate{1, 2, 3}.Equal(Coordinate{1, 2, 4}))
}
