package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBomberSpeed(t *testing.T) {
	b := newBomber()
	assert.Equal(t, int64(8800), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 7}, false, false))
	assert.Equal(t, int64(8800), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 0}, false, false))
	assert.Equal(t, int64(17000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, false, false))
	assert.Equal(t, int64(17000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, false, false))
	assert.Equal(t, int64(22000), b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}, false, true))
}
