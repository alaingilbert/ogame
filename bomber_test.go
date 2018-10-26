package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBomberSpeed(t *testing.T) {
	b := newBomber()
	assert.Equal(t, 8800, b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 7}))
	assert.Equal(t, 8800, b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 0}))
	assert.Equal(t, 17000, b.GetSpeed(Researches{ImpulseDrive: 6, HyperspaceDrive: 8}))
}
