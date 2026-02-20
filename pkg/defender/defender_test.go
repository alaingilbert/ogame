package defender

import (
	"testing"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := DefaultConfig()
	d := New(nil, cfg)
	assert.NotNil(t, d)
	assert.False(t, d.IsRunning())
	assert.Equal(t, cfg.CheckIntervalSec, d.GetConfig().CheckIntervalSec)
}

func TestGetSetConfig(t *testing.T) {
	cfg := DefaultConfig()
	d := New(nil, cfg)

	newCfg := cfg
	newCfg.CheckIntervalSec = 600
	newCfg.EvacSpeed = ogame.FiftyPercent
	newCfg.EvacSecondsBefore = 25
	newCfg.RecallSecondsAfter = 15
	newCfg.Origin = CelestialOrigin{Random: false, CelestialID: ogame.CelestialID(99999)}
	d.SetConfig(newCfg)

	got := d.GetConfig()
	assert.Equal(t, 600, got.CheckIntervalSec)
	assert.Equal(t, ogame.FiftyPercent, got.EvacSpeed)
	assert.Equal(t, 25, got.EvacSecondsBefore)
	assert.Equal(t, 15, got.RecallSecondsAfter)
	assert.False(t, got.Origin.Random)
	assert.Equal(t, ogame.CelestialID(99999), got.Origin.CelestialID)
}

func TestEvacuatedFleetTracking(t *testing.T) {
	cfg := DefaultConfig()
	d := New(nil, cfg)

	celestialID := ogame.CelestialID(12345)
	fleetID := ogame.FleetID(67890)

	d.mu.Lock()
	d.evacuatedFleets[celestialID] = evacuatedFleet{
		FleetID: fleetID,
		AttackEvent: ogame.AttackEvent{
			ArrivalTime: time.Now().Add(5 * time.Minute),
		},
		SentAt: time.Now(),
	}
	d.mu.Unlock()

	d.mu.Lock()
	evac, exists := d.evacuatedFleets[celestialID]
	d.mu.Unlock()

	assert.True(t, exists)
	assert.Equal(t, fleetID, evac.FleetID)
}

func TestStopWithoutStart(t *testing.T) {
	cfg := DefaultConfig()
	d := New(nil, cfg)
	// Should not panic.
	d.Stop()
	assert.False(t, d.IsRunning())
}

func TestStartDisabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Enabled = false
	d := New(nil, cfg)
	d.Start()
	assert.False(t, d.IsRunning())
}
