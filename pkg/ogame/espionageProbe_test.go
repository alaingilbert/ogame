package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEspionageProbeFuelConsumption(t *testing.T) {
	ep := newEspionageProbe()
	assert.Equal(t, int64(1), ep.GetFuelConsumption(Researches{}, 1, false))
	assert.Equal(t, int64(0), ep.GetFuelConsumption(Researches{}, 1, true))
	assert.Equal(t, int64(0), ep.GetFuelConsumption(Researches{}, 0.5, false))
}
