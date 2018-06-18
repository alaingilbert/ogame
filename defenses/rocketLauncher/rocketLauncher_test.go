package rocketLauncher

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestConstructionTime(t *testing.T) {
	rl := New()
	assert.Equal(t, 82, rl.ConstructionTime(1, 7, ogame.Facilities{Shipyard: 4}))
	assert.Equal(t, 164, rl.ConstructionTime(2, 7, ogame.Facilities{Shipyard: 4}))
}
