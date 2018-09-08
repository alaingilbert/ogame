package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCrystalMineConstructionTime(t *testing.T) {
	cm := newCrystalMine()
	assert.Equal(t, 75*time.Second, cm.ConstructionTime(5, 6, Facilities{}))
}
