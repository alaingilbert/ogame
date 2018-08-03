package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRocketLauncherConstructionTime(t *testing.T) {
	rl := NewRocketLauncher()
	assert.Equal(t, 82, rl.ConstructionTime(1, 7, Facilities{Shipyard: 4}))
	assert.Equal(t, 164, rl.ConstructionTime(2, 7, Facilities{Shipyard: 4}))
}
