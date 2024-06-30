package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathfinderSpeed(t *testing.T) {
	pf := newPathfinder()
	assert.Equal(t, int64(12000), pf.GetSpeed(Researches{}, LfBonuses{}, NoClass, NoAllianceClass))
	assert.Equal(t, int64(26400), pf.GetSpeed(Researches{HyperspaceDrive: 4}, LfBonuses{}, NoClass, NoAllianceClass))
}
