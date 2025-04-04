package v11_13_0

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGetConstructions(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.13.0/en/overview.html")
	clock := clockwork.NewFakeClockAt(time.Date(2024, 4, 11, 21, 24, 7, 0, time.UTC))
	constructions, _ := ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.DeuteriumTankID, constructions.Building.ID)
	assert.Equal(t, int64(243), int64(constructions.Building.Countdown.Seconds()))
	assert.Equal(t, int64(11), constructions.Building.Level)
	assert.Equal(t, ogame.IonTechnologyID, constructions.Research.ID)
	assert.Equal(t, int64(52), int64(constructions.Research.Countdown.Seconds()))
	assert.Equal(t, int64(6), constructions.Research.Level)
	assert.Equal(t, ogame.ResidentialSectorID, constructions.LfBuilding.ID)
	assert.Equal(t, int64(414), int64(constructions.LfBuilding.Countdown.Seconds()))
	assert.Equal(t, int64(42), constructions.LfBuilding.Level)
	assert.Equal(t, ogame.IntergalacticEnvoysID, constructions.LfResearch.ID)
	assert.Equal(t, int64(25972), int64(constructions.LfResearch.Countdown.Seconds()))
	assert.Equal(t, int64(15), constructions.LfResearch.Level)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.13.0/en/shipyard.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 3, len(prods))
	assert.Equal(t, int64(632), secs)
	assert.Equal(t, ogame.DestroyerID, prods[0].ID)
	assert.Equal(t, int64(28), prods[0].Nbr)
	assert.Equal(t, ogame.DestroyerID, prods[1].ID)
	assert.Equal(t, int64(30), prods[1].Nbr)
	assert.Equal(t, ogame.LightFighterID, prods[2].ID)
	assert.Equal(t, int64(1), prods[2].Nbr)
}

func TestExtractFleet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.13.0/en/movement.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets, _ := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(10495), fleets[0].ArriveIn)
	assert.Equal(t, int64(20995), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 103, Position: 11, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 105, Position: 11, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(8568846), fleets[0].ID)
	assert.Equal(t, int64(60), fleets[0].Ships.Destroyer)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}

func TestExtractFleet2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.13.0/en/movement2.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets, _ := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(3416), fleets[0].ArriveIn)
	assert.Equal(t, int64(3908), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 103, Position: 11, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 103, Position: 16, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Expedition, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(8573942), fleets[0].ID)
	assert.Equal(t, int64(1), fleets[0].Ships.LightFighter)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}
