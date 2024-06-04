package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSolarPlant_Production(t *testing.T) {
	sp := newSolarPlant()
	assert.Equal(t, int64(9200), sp.Production(29))
}

func TestSolarPlant_ConstructionTime(t *testing.T) {
	sp := newSolarPlant()
	ct := sp.ConstructionTime(1, 6, Facilities{RoboticsFactory: 10, NaniteFactory: 7}, LfBonuses{}, NoClass, false)
	assert.Equal(t, time.Second, ct)
}
