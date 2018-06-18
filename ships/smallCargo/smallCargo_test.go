package smallCargo

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestConstructionTime(t *testing.T) {
	sc := New()
	assert.Equal(t, 164, sc.ConstructionTime(1, 7, ogame.Facilities{Shipyard: 4}))
	assert.Equal(t, 328, sc.ConstructionTime(2, 7, ogame.Facilities{Shipyard: 4}))
}
