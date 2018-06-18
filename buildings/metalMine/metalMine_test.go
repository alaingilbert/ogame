package metalMine

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestProduction(t *testing.T) {
	mm := New()
	assert.Equal(t, 30, mm.Production(1, 1, 0))
	assert.Equal(t, 63, mm.Production(1, 1, 1))
	assert.Equal(t, 120, mm.Production(4, 1, 0))
	assert.Equal(t, 252, mm.Production(4, 1, 1))
}

func TestConstructionTime(t *testing.T) {
	mm := New()
	ct := mm.ConstructionTime(20, 7, ogame.Facilities{RoboticsFactory: 3})
	assert.Equal(t, 8550, ct)
}
