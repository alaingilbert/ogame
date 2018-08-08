package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombustionDriveCost(t *testing.T) {
	cd := newCombustionDrive()
	assert.Equal(t, Resources{Metal: 12800, Deuterium: 19200}, cd.GetPrice(6))
}
