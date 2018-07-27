package fusionReactor

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCapacity(t *testing.T) {
	fr := New()
	assert.Equal(t, 38, fr.ConstructionTime(2, 7, ogame.Facilities{RoboticsFactory: 3}))
}
