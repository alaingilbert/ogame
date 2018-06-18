package armourTechnology

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestIsAvailable(t *testing.T) {
	resourcesBuildings := ogame.ResourcesBuildings{}
	facilities := ogame.Facilities{ResearchLab: 2}
	researches := ogame.Researches{}
	b := New()
	avail := b.IsAvailable(resourcesBuildings, facilities, researches, 0)
	assert.True(t, avail)
}

func TestIsAvailable_NoBuilding(t *testing.T) {
	resourcesBuildings := ogame.ResourcesBuildings{}
	facilities := ogame.Facilities{ResearchLab: 1}
	researches := ogame.Researches{}
	b := New()
	avail := b.IsAvailable(resourcesBuildings, facilities, researches, 0)
	assert.False(t, avail)
}
