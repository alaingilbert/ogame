package ogame

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemInfos_Position(t *testing.T) {
	si := SystemInfos{}
	var nilPlanetInfo *PlanetInfos
	assert.Equal(t, nilPlanetInfo, si.Position(0))
	assert.Equal(t, nilPlanetInfo, si.Position(16))
}

func TestSystemInfos_Each(t *testing.T) {
	si := SystemInfos{}
	i := 0
	si.Each(func(pi *PlanetInfos) {
		i++
	})
	assert.Equal(t, len(si.planets), i)
}

func TestSystemInfos_MarshalJSON(t *testing.T) {
	planetInfos := PlanetInfos{
		ID:         1,
		Activity:   15,
		Name:       "name",
		Img:        "img",
		Coordinate: Coordinate{1, 2, 3, PlanetType},
	}
	planetInfos.Debris.Metal = 1
	planetInfos.Debris.Crystal = 2
	planetInfos.Debris.RecyclersNeeded = 3
	planetInfos.Player.ID = 1
	planetInfos.Player.Name = "player name"
	planetInfos.Player.Rank = 2
	si := SystemInfos{}
	si.galaxy = 1
	si.system = 2
	si.planets[1] = &planetInfos
	by, _ := json.Marshal(si)
	expected := `{"Galaxy":1,"System":2,` +
		`"Planets":[null,` +
		`{"ID":1,"Activity":15,"Name":"name","Img":"img","Coordinate":{"Galaxy":1,"System":2,"Position":3,"Type":1},` +
		`"Administrator":false,"Inactive":false,"Vacation":false,"StrongPlayer":false,"Newbie":false,` +
		`"HonorableTarget":false,"Banned":false,"Debris":{"Metal":1,"Crystal":2,"RecyclersNeeded":3},"Moon":null,` +
		`"Player":{"ID":1,"Name":"player name","Rank":2,"IsBandit":false,"IsStarlord":false},"Alliance":null},` +
		`null,null,null,null,null,null,null,null,null,null,null,null,null],"ExpeditionDebris":{"Metal":0,"Crystal":0,"PathfindersNeeded":0}}`
	assert.Equal(t, expected, string(by))
}
