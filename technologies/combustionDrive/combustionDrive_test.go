package combustionDrive

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	cd := New()
	assert.Equal(t, ogame.Resources{Metal: 12800, Deuterium: 19200}, cd.GetPrice(6))
}
