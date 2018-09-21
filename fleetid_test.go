package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFleetID_String(t *testing.T) {
	assert.Equal(t, "12345", FleetID(12345).String())
}
