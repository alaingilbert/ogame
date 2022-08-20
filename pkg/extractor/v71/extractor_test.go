package v71

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestExtractAttacksACSAttackSelf(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v8.6/en/eventlist_acs_attack_self.html")
	ownCoords := []ogame.Coordinate{{4, 116, 9, ogame.PlanetType}}
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(1), attacks[0].Ships.LightFighter)
}

func TestExtractAttacksACS_v71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v7.1/en/eventlist_acs.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(200), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(9), attacks[0].Ships.SmallCargo)
}

func TestExtractAttacksACS_v72(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v7.2/en/eventlist_multipleACS.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 3, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(14028), attacks[0].ID)
	assert.Equal(t, int64(14029), attacks[1].ID)
	assert.Equal(t, int64(673019), attacks[2].ID)
}
