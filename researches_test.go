package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResearchesByID(t *testing.T) {
	r := Researches{
		EnergyTechnology:             1,
		LaserTechnology:              2,
		IonTechnology:                3,
		HyperspaceTechnology:         4,
		PlasmaTechnology:             5,
		CombustionDrive:              6,
		ImpulseDrive:                 7,
		HyperspaceDrive:              8,
		EspionageTechnology:          9,
		ComputerTechnology:           10,
		Astrophysics:                 11,
		IntergalacticResearchNetwork: 12,
		GravitonTechnology:           13,
		WeaponsTechnology:            14,
		ShieldingTechnology:          15,
		ArmourTechnology:             16,
	}
	assert.Equal(t, int64(0), r.ByID(123456))
	assert.Equal(t, int64(1), r.ByID(EnergyTechnologyID))
	assert.Equal(t, int64(2), r.ByID(LaserTechnologyID))
	assert.Equal(t, int64(3), r.ByID(IonTechnologyID))
	assert.Equal(t, int64(4), r.ByID(HyperspaceTechnologyID))
	assert.Equal(t, int64(5), r.ByID(PlasmaTechnologyID))
	assert.Equal(t, int64(6), r.ByID(CombustionDriveID))
	assert.Equal(t, int64(7), r.ByID(ImpulseDriveID))
	assert.Equal(t, int64(8), r.ByID(HyperspaceDriveID))
	assert.Equal(t, int64(9), r.ByID(EspionageTechnologyID))
	assert.Equal(t, int64(10), r.ByID(ComputerTechnologyID))
	assert.Equal(t, int64(11), r.ByID(AstrophysicsID))
	assert.Equal(t, int64(12), r.ByID(IntergalacticResearchNetworkID))
	assert.Equal(t, int64(13), r.ByID(GravitonTechnologyID))
	assert.Equal(t, int64(14), r.ByID(WeaponsTechnologyID))
	assert.Equal(t, int64(15), r.ByID(ShieldingTechnologyID))
	assert.Equal(t, int64(16), r.ByID(ArmourTechnologyID))
}

func TestResearchesString(t *testing.T) {
	r := Researches{
		EnergyTechnology:             1,
		LaserTechnology:              2,
		IonTechnology:                3,
		HyperspaceTechnology:         4,
		PlasmaTechnology:             5,
		CombustionDrive:              6,
		ImpulseDrive:                 7,
		HyperspaceDrive:              8,
		EspionageTechnology:          9,
		ComputerTechnology:           10,
		Astrophysics:                 11,
		IntergalacticResearchNetwork: 12,
		GravitonTechnology:           13,
		WeaponsTechnology:            14,
		ShieldingTechnology:          15,
		ArmourTechnology:             16,
	}
	expected := "\n" +
		"             Energy Technology: 1\n" +
		"              Laser Technology: 2\n" +
		"                Ion Technology: 3\n" +
		"         Hyperspace Technology: 4\n" +
		"             Plasma Technology: 5\n" +
		"              Combustion Drive: 6\n" +
		"                 Impulse Drive: 7\n" +
		"              Hyperspace Drive: 8\n" +
		"          Espionage Technology: 9\n" +
		"           Computer Technology: 10\n" +
		"                  Astrophysics: 11\n" +
		"Intergalactic Research Network: 12\n" +
		"           Graviton Technology: 13\n" +
		"            Weapons Technology: 14\n" +
		"          Shielding Technology: 15\n" +
		"             Armour Technology: 16"
	assert.Equal(t, expected, r.String())
}
