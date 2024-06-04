package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEspionageProbeFuelConsumption(t *testing.T) {
	ep := newEspionageProbe()
	assert.Equal(t, int64(1), ep.GetFuelConsumption(Researches{}, LfBonuses{}, NoClass, 1))
	assert.Equal(t, int64(0), ep.GetFuelConsumption(Researches{}, LfBonuses{}, General, 1))
	assert.Equal(t, int64(0), ep.GetFuelConsumption(Researches{}, LfBonuses{}, NoClass, 0.5))
}

func TestEspionageProbe_GetCargoCapacity(t *testing.T) {
	ep := newEspionageProbe()
	assert.Equal(t, int64(8), ep.GetCargoCapacity(Researches{HyperspaceTechnology: 14}, LfBonuses{}, NoClass, 0.05, true))
}
