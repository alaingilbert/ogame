package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByID(t *testing.T) {
	assert.Equal(t, int64(0), ShipsInfos{}.ByID(123456))
}

func TestShipsInfos_Cargo(t *testing.T) {
	ships := ShipsInfos{
		SmallCargo: 2,
		LargeCargo: 2,
	}
	techs := Researches{}
	assert.Equal(t, int64(60000), ships.Cargo(techs, false, false, false))
}

func TestShipsInfos_FleetValue(t *testing.T) {
	ships := ShipsInfos{
		SmallCargo: 2,
		LargeCargo: 2,
	}
	assert.Equal(t, int64(32000), ships.FleetValue())
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
	assert.Equal(t, int64(2), s1.LightFighter)
	assert.Equal(t, int64(4), s1.HeavyFighter)
}

func TestSet(t *testing.T) {
	s := ShipsInfos{}
	s.Set(BattleshipID, 1)
	s.Set(DeathstarID, 2)
	s.Set(SolarSatelliteID, 4)
	assert.Equal(t, int64(1), s.ByID(BattleshipID))
	assert.Equal(t, int64(2), s.ByID(DeathstarID))
	assert.Equal(t, int64(4), s.ByID(SolarSatelliteID))
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
		Crawler:        15,
		Reaper:         16,
		Pathfinder:     17,
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
		"Solar Satellite: 14\n" +
		"        Crawler: 15\n" +
		"         Reaper: 16\n" +
		"     Pathfinder: 17"
	assert.Equal(t, expected, s.String())
}

func TestShipsInfos_Equal(t *testing.T) {
	ships := ShipsInfos{SmallCargo: 2, LargeCargo: 3}
	assert.True(t, ships.Equal(ShipsInfos{SmallCargo: 2, LargeCargo: 3}))
	assert.False(t, ships.Equal(ShipsInfos{SmallCargo: 2, LargeCargo: 3, EspionageProbe: 4}))
}

func TestShipsInfos_FleetCost(t *testing.T) {
	assert.Equal(t, Resources{Metal: 22000, Crystal: 22000}, ShipsInfos{SmallCargo: 2, LargeCargo: 3}.FleetCost())
}

func TestShipsInfos_CountShips(t *testing.T) {
	assert.Equal(t, int64(5), ShipsInfos{SmallCargo: 2, LargeCargo: 3}.CountShips())
}

func TestShipsInfos_Has(t *testing.T) {
	ships := ShipsInfos{SmallCargo: 2, LargeCargo: 3}
	assert.True(t, ships.Has(ShipsInfos{SmallCargo: 1, LargeCargo: 2}))
	assert.True(t, ships.Has(ShipsInfos{SmallCargo: 2, LargeCargo: 3}))
	assert.False(t, ships.Has(ShipsInfos{SmallCargo: 2, LargeCargo: 4}))
	assert.True(t, ships.Has(ShipsInfos{SmallCargo: 2}))
}

func TestShipsInfos_HasShips(t *testing.T) {
	assert.True(t, ShipsInfos{SmallCargo: 2, LargeCargo: 3}.HasShips())
	assert.False(t, ShipsInfos{}.HasShips())
	assert.True(t, ShipsInfos{SolarSatellite: 1}.HasShips())
}

func TestShipsInfos_ToQuantifiables(t *testing.T) {
	assert.Equal(t, []Quantifiable{{SmallCargoID, 1}, {LargeCargoID, 2}}, ShipsInfos{SmallCargo: 1, LargeCargo: 2}.ToQuantifiables())
}

func TestShipsInfos_FromQuantifiables(t *testing.T) {
	assert.Equal(t, ShipsInfos{SmallCargo: 1, LargeCargo: 2}, ShipsInfos{}.FromQuantifiables([]Quantifiable{{SmallCargoID, 1}, {LargeCargoID, 2}}))
}

func TestShipsInfos_Speed(t *testing.T) {
	assert.Equal(t, int64(20250), ShipsInfos{LargeCargo: 2}.Speed(Researches{CombustionDrive: 17}, false, false))
	assert.Equal(t, int64(20250), ShipsInfos{LargeCargo: 2, SolarSatellite: 1}.Speed(Researches{CombustionDrive: 17}, false, false))
}

func TestShipsInfos_ToPtr(t *testing.T) {
	ships := ShipsInfos{SmallCargo: 2, LargeCargo: 3}
	shipsPtr := ships.ToPtr()
	assert.Equal(t, &ships, shipsPtr)
}
