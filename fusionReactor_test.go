package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFusionReactorCapacity(t *testing.T) {
	fr := newFusionReactor()
	assert.Equal(t, 38*time.Second, fr.ConstructionTime(2, 7, Facilities{RoboticsFactory: 3}))
}
