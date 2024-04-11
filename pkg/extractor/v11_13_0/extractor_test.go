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
