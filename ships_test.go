package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByID(t *testing.T) {
	assert.Equal(t, 0, ShipsInfos{}.ByID(123456))
}

func TestShipsInfos_Cargo(t *testing.T) {
	ships := ShipsInfos{
		SmallCargo: 2,
		LargeCargo: 2,
	}
	techs := Researches{}
	assert.Equal(t, 60000, ships.Cargo(techs))
}

func TestShipsInfos_Add(t *testing.T) {
	s1 := ShipsInfos{
		LightFighter:   1,
		HeavyFighter:   2,
		Cruiser:        3,
		Battleship:     4,
		Battlecruiser:  5,
		Bomber:         6,
		Destroyer:      7,
		Deathstar:      8,
		SmallCargo:     9,
		LargeCargo:     10,
		ColonyShip:     11,
		Recycler:       12,
		EspionageProbe: 13,
		SolarSatellite: 14,
	}
	s1.Add(ShipsInfos{
		LightFighter:   1,
		HeavyFighter:   2,
		Cruiser:        3,
		Battleship:     4,
		Battlecruiser:  5,
		Bomber:         6,
		Destroyer:      7,
		Deathstar:      8,
		SmallCargo:     9,
		LargeCargo:     10,
		ColonyShip:     11,
		Recycler:       12,
		EspionageProbe: 13,
		SolarSatellite: 14,
	})
	assert.Equal(t, 2, s1.LightFighter)
	assert.Equal(t, 4, s1.HeavyFighter)
}

func TestSet(t *testing.T) {
	s := ShipsInfos{}
	s.Set(BattleshipID, 1)
	s.Set(DeathstarID, 2)
	s.Set(SolarSatelliteID, 4)
	assert.Equal(t, 1, s.ByID(BattleshipID))
	assert.Equal(t, 2, s.ByID(DeathstarID))
	assert.Equal(t, 4, s.ByID(SolarSatelliteID))
}

func TestShipsInfos_String(t *testing.T) {
	s := ShipsInfos{
		LightFighter:   1,
		HeavyFighter:   2,
		Cruiser:        3,
		Battleship:     4,
		Battlecruiser:  5,
		Bomber:         6,
		Destroyer:      7,
		Deathstar:      8,
		SmallCargo:     9,
		LargeCargo:     10,
		ColonyShip:     11,
		Recycler:       12,
		EspionageProbe: 13,
		SolarSatellite: 14,
	}
	expected := "\n" +
		"  Light Fighter: 1\n" +
		"  Heavy Fighter: 2\n" +
		"        Cruiser: 3\n" +
		"     Battleship: 4\n" +
		"  Battlecruiser: 5\n" +
		"         Bomber: 6\n" +
		"      Destroyer: 7\n" +
		"      Deathstar: 8\n" +
		"    Small Cargo: 9\n" +
		"    Large Cargo: 10\n" +
		"    Colony Ship: 11\n" +
		"       Recycler: 12\n" +
		"Espionage Probe: 13\n" +
		"Solar Satellite: 14"
	assert.Equal(t, expected, s.String())
}
