package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBomberSpeed(t *testing.T) {
	b := newBomber()
	assert.Equal(t, int64(8800), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 7}, LfBonuses{}, NoClass, NoAllianceClass))
	assert.Equal(t, int64(8800), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 0}, LfBonuses{}, NoClass, NoAllianceClass))
	assert.Equal(t, int64(17000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, LfBonuses{}, NoClass, NoAllianceClass))
	assert.Equal(t, int64(17000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, LfBonuses{}, NoClass, NoAllianceClass))
	assert.Equal(t, int64(22000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, LfBonuses{}, General, NoAllianceClass))

	lfBonuses := LfBonuses{LfShipBonuses: make(LfShipBonuses)}
	lfBonuses.LfShipBonuses[BomberID] = LfShipBonus{Speed: 0.2478}
	assert.Equal(t, int64(28739), b.GetSpeed(Researches{ImpulseDrive: 15, HyperspaceDrive: 15}, lfBonuses, Discoverer, NoAllianceClass))
}
