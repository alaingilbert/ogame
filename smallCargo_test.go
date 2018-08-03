package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmallCargoConstructionTime(t *testing.T) {
	sc := NewSmallCargo()
	assert.Equal(t, 164, sc.ConstructionTime(1, 7, Facilities{Shipyard: 4}))
	assert.Equal(t, 328, sc.ConstructionTime(2, 7, Facilities{Shipyard: 4}))
}
