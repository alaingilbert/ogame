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
	buildingID, buildingCountdown, researchID, researchCountdown, lfBuildingID, lfBuildingCountdown, lfResearchID, lfResearchCountdown := ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.DeuteriumTankID, buildingID)
	assert.Equal(t, int64(243), buildingCountdown)
	assert.Equal(t, ogame.IonTechnologyID, researchID)
	assert.Equal(t, int64(52), researchCountdown)
	assert.Equal(t, ogame.ResidentialSectorID, lfBuildingID)
	assert.Equal(t, int64(414), lfBuildingCountdown)
	assert.Equal(t, ogame.IntergalacticEnvoysID, lfResearchID)
	assert.Equal(t, int64(25972), lfResearchCountdown)
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
