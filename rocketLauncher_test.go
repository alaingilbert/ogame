package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRocketLauncherConstructionTime(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, 82*time.Second, rl.ConstructionTime(1, 7, Facilities{Shipyard: 4}))
	assert.Equal(t, 164*time.Second, rl.ConstructionTime(2, 7, Facilities{Shipyard: 4}))
}
