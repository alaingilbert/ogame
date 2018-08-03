package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImpulseDriveCost(t *testing.T) {
	id := NewImpulseDrive()
	assert.Equal(t, Resources{Metal: 8000, Crystal: 16000, Deuterium: 2400}, id.GetPrice(3))
}
