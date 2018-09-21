package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolarPlant_Production(t *testing.T) {
	sp := newSolarPlant()
	assert.Equal(t, 9200, sp.Production(29))
}
