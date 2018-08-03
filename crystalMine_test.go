package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrystalMineConstructionTime(t *testing.T) {
	cm := NewCrystalMine()
	assert.Equal(t, 75, cm.ConstructionTime(5, 6, Facilities{}))
}
