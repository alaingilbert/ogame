package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNaniteFactoryConstructionTime(t *testing.T) {
	nf := newNaniteFactory()
	assert.Equal(t, 28051*time.Second, nf.ConstructionTime(2, 7, Facilities{RoboticsFactory: 10, NaniteFactory: 1}, false, true))
	assert.Equal(t, 28051*time.Second, nf.ConstructionTime(3, 7, Facilities{RoboticsFactory: 10, NaniteFactory: 2}, false, true))
	assert.Equal(t, 22040*time.Second, nf.ConstructionTime(6, 7, Facilities{RoboticsFactory: 13, NaniteFactory: 5}, false, true))

	assert.Equal(t, 39272*time.Second, nf.ConstructionTime(1, 5, Facilities{RoboticsFactory: 10, NaniteFactory: 0}, false, false))
	assert.Equal(t, 39272*time.Second, nf.ConstructionTime(2, 5, Facilities{RoboticsFactory: 10, NaniteFactory: 1}, false, false))
	assert.Equal(t, 39272*time.Second, nf.ConstructionTime(3, 5, Facilities{RoboticsFactory: 10, NaniteFactory: 2}, false, false))
}
