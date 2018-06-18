package impulseDrive

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	id := New()
	assert.Equal(t, ogame.Resources{Metal: 8000, Crystal: 16000, Deuterium: 2400}, id.GetPrice(3))
}
