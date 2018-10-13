package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants_DestinationType_String(t *testing.T) {
	assert.Equal(t, "planet", PlanetDest.String())
	assert.Equal(t, "moon", MoonDest.String())
	assert.Equal(t, "debris", DebrisDest.String())
	assert.Equal(t, "123", DestinationType(123).String())
}

func TestConstants_Speed_String(t *testing.T) {
	assert.Equal(t, "10%", Speed(1).String())
	assert.Equal(t, "20%", Speed(2).String())
	assert.Equal(t, "30%", Speed(3).String())
	assert.Equal(t, "40%", Speed(4).String())
	assert.Equal(t, "50%", Speed(5).String())
	assert.Equal(t, "60%", Speed(6).String())
	assert.Equal(t, "70%", Speed(7).String())
	assert.Equal(t, "80%", Speed(8).String())
	assert.Equal(t, "90%", Speed(9).String())
	assert.Equal(t, "100%", Speed(10).String())
	assert.Equal(t, "11", Speed(11).String())
}

func TestConstants_MissionID_String(t *testing.T) {
	assert.Equal(t, "Attack", MissionID(1).String())
	assert.Equal(t, "GroupedAttack", MissionID(2).String())
	assert.Equal(t, "Transport", MissionID(3).String())
	assert.Equal(t, "Park", MissionID(4).String())
	assert.Equal(t, "ParkInThatAlly", MissionID(5).String())
	assert.Equal(t, "Spy", MissionID(6).String())
	assert.Equal(t, "Colonize", MissionID(7).String())
	assert.Equal(t, "RecycleDebrisField", MissionID(8).String())
	assert.Equal(t, "Destroy", MissionID(9).String())
	assert.Equal(t, "MissileAttack", MissionID(10).String())
	assert.Equal(t, "Expedition", MissionID(15).String())
	assert.Equal(t, "16", MissionID(16).String())
}
