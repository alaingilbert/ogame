package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceSettings_String(t *testing.T) {
	r := ResourceSettings{
		MetalMine:            1,
		CrystalMine:          2,
		DeuteriumSynthesizer: 3,
		SolarPlant:           4,
		FusionReactor:        5,
		SolarSatellite:       6,
	}
	expected := "\n" +
		"           Metal Mine: 1\n" +
		"         Crystal Mine: 2\n" +
		"Deuterium Synthesizer: 3\n" +
		"          Solar Plant: 4\n" +
		"       Fusion Reactor: 5\n" +
		"      Solar Satellite: 6"
	assert.Equal(t, expected, r.String())
}
