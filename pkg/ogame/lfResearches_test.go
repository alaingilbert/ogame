package ogame

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntergalacticEnvoysConstructionTime(t *testing.T) {
	ie := newIntergalacticEnvoys()
	assert.Equal(t, (6*60+0)*time.Second, ie.ConstructionTime(2, 8, Facilities{}, false, false))
}
