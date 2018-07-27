package crystalMine

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestConstructionTime(t *testing.T) {
	cm := New()
	assert.Equal(t, 75, cm.ConstructionTime(5, 6, ogame.Facilities{}))
}
