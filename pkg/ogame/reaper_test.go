package ogame

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReaperSpeed(t *testing.T) {
	r := newReaper()
	lfBonuses := LfBonuses{LfShipBonuses: make(LfShipBonuses)}
	lfBonuses.LfShipBonuses[ReaperID] = LfShipBonus{Speed: 0.2478}
	assert.Equal(t, int64(40235), r.GetSpeed(Researches{HyperspaceDrive: 15}, lfBonuses, Discoverer))
}
