package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSolarSatelliteSpeed(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, int64(0), ss.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}, false))
}

func TestSolarSatellite_GetLevel(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, int64(0), ss.GetLevel(ResourcesBuildings{SolarSatellite: 10}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
}

func TestSolarSatellite_Production(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, int64(1326), ss.Production(Temperature{-23, 17}, 51))
}

func TestSolarSatellite_ConstructionTime(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, 1*time.Second, ss.ConstructionTime(1, 7, Facilities{Shipyard: 12, NaniteFactory: 6, RoboticsFactory: 10}))
	assert.Equal(t, 6*time.Second, ss.ConstructionTime(1, 7, Facilities{Shipyard: 1, NaniteFactory: 5, RoboticsFactory: 10}))
	assert.Equal(t, 102*time.Second, ss.ConstructionTime(1, 7, Facilities{Shipyard: 3, NaniteFactory: 0, RoboticsFactory: 10}))
}
