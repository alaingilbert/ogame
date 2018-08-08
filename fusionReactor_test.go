package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFusionReactorCapacity(t *testing.T) {
	fr := newFusionReactor()
	assert.Equal(t, 38, fr.ConstructionTime(2, 7, Facilities{RoboticsFactory: 3}))
}
