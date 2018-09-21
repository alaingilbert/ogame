package ogame

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields_HasFieldAvailable(t *testing.T) {
	assert.True(t, Fields{Built: 10, Total: 11}.HasFieldAvailable())
	assert.False(t, Fields{Built: 11, Total: 11}.HasFieldAvailable())
}

func TestTemperature_Mean(t *testing.T) {
	assert.Equal(t, 5, Temperature{Min: 0, Max: 10}.Mean())
	assert.Equal(t, 0, Temperature{Min: -10, Max: 10}.Mean())
}

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
