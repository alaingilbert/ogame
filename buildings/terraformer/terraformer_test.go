package terraformer

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestIsAvailable(t *testing.T) {
	resourcesBuildings := ogame.ResourcesBuildings{}
	facilities := ogame.Facilities{NaniteFactory: 1}
	researches := ogame.Researches{EnergyTechnology: 12}
	b := New()
	avail := b.IsAvailable(resourcesBuildings, facilities, researches, 0)
	assert.True(t, avail)
}

func TestIsAvailable_NoTech(t *testing.T) {
	resourcesBuildings := ogame.ResourcesBuildings{}
	facilities := ogame.Facilities{NaniteFactory: 1}
	researches := ogame.Researches{EnergyTechnology: 11}
	b := New()
	avail := b.IsAvailable(resourcesBuildings, facilities, researches, 0)
	assert.False(t, avail)
}
