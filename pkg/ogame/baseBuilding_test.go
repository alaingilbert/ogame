package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseBuilding_GetLevel(t *testing.T) {
	var bb struct {
		BaseBuilding
	}
	bb.ID = ID(123456)
	assert.Equal(t, int64(0), bb.GetLevel(ResourcesBuildings{}, Facilities{}, Researches{}))
}
