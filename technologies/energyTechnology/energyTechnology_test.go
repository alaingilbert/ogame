package energyTechnology

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestConstructionTime(t *testing.T) {
	mm := New()
	ct := mm.ConstructionTime(5, 7, ogame.Facilities{ResearchLab: 3})
	assert.Equal(t, 1645, ct)
}
