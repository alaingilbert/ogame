package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecyclerSpeed(t *testing.T) {
	r := newRecycler()
	assert.Equal(t, int64(3200), r.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 1, HyperspaceDrive: 1}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(17600), r.GetSpeed(Researches{CombustionDrive: 1, ImpulseDrive: 17, HyperspaceDrive: 10}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(33000), r.GetSpeed(Researches{CombustionDrive: 1, ImpulseDrive: 17, HyperspaceDrive: 15}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(18400), r.GetSpeed(Researches{CombustionDrive: 1, ImpulseDrive: 18, HyperspaceDrive: 10}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(34800), r.GetSpeed(Researches{CombustionDrive: 1, ImpulseDrive: 17, HyperspaceDrive: 16}, LfBonuses{}, NoClass))
}
