package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetalMineProduction(t *testing.T) {
	mm := NewMetalMine()
	assert.Equal(t, 30, mm.Production(1, 1, 0))
	assert.Equal(t, 63, mm.Production(1, 1, 1))
	assert.Equal(t, 120, mm.Production(4, 1, 0))
	assert.Equal(t, 252, mm.Production(4, 1, 1))
}

func TestMetalMineConstructionTime(t *testing.T) {
	mm := NewMetalMine()
	assert.Equal(t, 8550, mm.ConstructionTime(20, 7, Facilities{RoboticsFactory: 3}))
	assert.Equal(t, 30, mm.ConstructionTime(4, 6, Facilities{}))
}
