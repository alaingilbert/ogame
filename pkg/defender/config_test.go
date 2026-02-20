package defender

import (
	"testing"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.False(t, cfg.Enabled)
	assert.Equal(t, 300, cfg.CheckIntervalSec)
	assert.Equal(t, ogame.HundredPercent, cfg.EvacSpeed)
	assert.Equal(t, 30, cfg.EvacSecondsBefore)
	assert.Equal(t, 0, cfg.RecallSecondsAfter)
	assert.True(t, cfg.Origin.Random)
	assert.False(t, cfg.DumpResourcesToDefense)
	assert.False(t, cfg.DumpResourcesToDefenseMoons)
}

func TestValidate_MinCheckInterval(t *testing.T) {
	cfg := Config{CheckIntervalSec: 10, EvacSecondsBefore: 1}
	cfg.Validate()
	assert.Equal(t, 60, cfg.CheckIntervalSec)
	assert.Equal(t, 5, cfg.EvacSecondsBefore)
}

func TestValidate_NegativeRecall(t *testing.T) {
	cfg := Config{
		CheckIntervalSec:   300,
		EvacSecondsBefore:  20,
		RecallSecondsAfter: -5,
	}
	cfg.Validate()
	assert.Equal(t, 0, cfg.RecallSecondsAfter)
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := Config{
		CheckIntervalSec:   600,
		EvacSecondsBefore:  25,
		RecallSecondsAfter: 15,
	}
	cfg.Validate()
	assert.Equal(t, 600, cfg.CheckIntervalSec)
	assert.Equal(t, 25, cfg.EvacSecondsBefore)
	assert.Equal(t, 15, cfg.RecallSecondsAfter)
}

func TestCelestialOrigin(t *testing.T) {
	// Random origin
	origin := CelestialOrigin{Random: true}
	assert.True(t, origin.Random)
	assert.Equal(t, ogame.CelestialID(0), origin.CelestialID)

	// Specific origin
	origin = CelestialOrigin{Random: false, CelestialID: ogame.CelestialID(12345)}
	assert.False(t, origin.Random)
	assert.Equal(t, ogame.CelestialID(12345), origin.CelestialID)
}
