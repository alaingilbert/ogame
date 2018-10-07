package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoonFacilities_ByID(t *testing.T) {
	f := MoonFacilities{
		RoboticsFactory: 1,
		Shipyard:        2,
		LunarBase:       3,
		SensorPhalanx:   4,
		JumpGate:        5,
	}
	assert.Equal(t, 1, f.ByID(RoboticsFactoryID))
	assert.Equal(t, 2, f.ByID(ShipyardID))
	assert.Equal(t, 3, f.ByID(LunarBaseID))
	assert.Equal(t, 4, f.ByID(SensorPhalanxID))
	assert.Equal(t, 5, f.ByID(JumpGateID))
	assert.Equal(t, 0, f.ByID(ID(12345)))
}

func TestMoonFacilities_String(t *testing.T) {
	f := MoonFacilities{
		RoboticsFactory: 1,
		Shipyard:        2,
		LunarBase:       3,
		SensorPhalanx:   4,
		JumpGate:        5,
	}
	expected := "\n" +
		"RoboticsFactory: 1\n" +
		"       Shipyard: 2\n" +
		"     Lunar Base: 3\n" +
		" Sensor Phalanx: 4\n" +
		"      Jump Gate: 5"
	assert.Equal(t, expected, f.String())
}

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
	}
	assert.Equal(t, 1, f.ByID(RoboticsFactoryID))
	assert.Equal(t, 2, f.ByID(ShipyardID))
	assert.Equal(t, 3, f.ByID(ResearchLabID))
	assert.Equal(t, 4, f.ByID(AllianceDepotID))
	assert.Equal(t, 5, f.ByID(MissileSiloID))
	assert.Equal(t, 6, f.ByID(NaniteFactoryID))
	assert.Equal(t, 7, f.ByID(TerraformerID))
	assert.Equal(t, 8, f.ByID(SpaceDockID))
	assert.Equal(t, 0, f.ByID(ID(12345)))
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
	}
	expected := "\n" +
		"RoboticsFactory: 1\n" +
		"       Shipyard: 2\n" +
		"   Research Lab: 3\n" +
		" Alliance Depot: 4\n" +
		"   Missile Silo: 5\n" +
		" Nanite Factory: 6\n" +
		"    Terraformer: 7\n" +
		"     Space Dock: 8"
	assert.Equal(t, expected, f.String())
}
