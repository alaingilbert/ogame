package ogame

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDefense(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/defence.html")
	defense, _ := extractDefense(string(pageHTMLBytes))
	assert.Equal(t, 1, defense.RocketLauncher)
	assert.Equal(t, 2, defense.LightLaser)
	assert.Equal(t, 3, defense.HeavyLaser)
	assert.Equal(t, 4, defense.GaussCannon)
	assert.Equal(t, 5, defense.IonCannon)
	assert.Equal(t, 6, defense.PlasmaTurret)
	assert.Equal(t, 0, defense.SmallShieldDome)
	assert.Equal(t, 0, defense.LargeShieldDome)
	assert.Equal(t, 7, defense.AntiBallisticMissiles)
	assert.Equal(t, 8, defense.InterplanetaryMissiles)
}
