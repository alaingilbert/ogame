package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacilities_ByID(t *testing.T) {
	f := Facilities{
		RoboticsFactory: 1,
		Shipyard:        2,
		ResearchLab:     3,
		AllianceDepot:   4,
		MissileSilo:     5,
		NaniteFactory:   6,
		Terraformer:     7,
		SpaceDock:       8,
		LunarBase:       9,
		SensorPhalanx:   10,
		JumpGate:        11,
	}
	assert.Equal(t, int64(1), f.ByID(RoboticsFactoryID))
	assert.Equal(t, int64(2), f.ByID(ShipyardID))
	assert.Equal(t, int64(3), f.ByID(ResearchLabID))
	assert.Equal(t, int64(4), f.ByID(AllianceDepotID))
	assert.Equal(t, int64(5), f.ByID(MissileSiloID))
	assert.Equal(t, int64(6), f.ByID(NaniteFactoryID))
	assert.Equal(t, int64(7), f.ByID(TerraformerID))
	assert.Equal(t, int64(8), f.ByID(SpaceDockID))
	assert.Equal(t, int64(9), f.ByID(LunarBaseID))
	assert.Equal(t, int64(10), f.ByID(SensorPhalanxID))
	assert.Equal(t, int64(11), f.ByID(JumpGateID))
	assert.Equal(t, int64(0), f.ByID(ID(12345)))
}

func TestFacilities_String(t *testing.T) {
	f := Facilities{
		RoboticsFactory: 1,
		Shipyard:        2,
		ResearchLab:     3,
		AllianceDepot:   4,
		MissileSilo:     5,
		NaniteFactory:   6,
		Terraformer:     7,
		SpaceDock:       8,
		LunarBase:       9,
		SensorPhalanx:   10,
		JumpGate:        11,
	}
	expected := "\n" +
		"RoboticsFactory: 1\n" +
		"       Shipyard: 2\n" +
		"   Research Lab: 3\n" +
		" Alliance Depot: 4\n" +
		"   Missile Silo: 5\n" +
		" Nanite Factory: 6\n" +
		"    Terraformer: 7\n" +
		"     Space Dock: 8\n" +
		"     Lunar Base: 9\n" +
		" Sensor Phalanx: 10\n" +
		"      Jump Gate: 11"
	assert.Equal(t, expected, f.String())
}
