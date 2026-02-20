package defender

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

// Defender monitors for incoming attacks and evacuates fleets + resources.
type Defender struct {
	bot    wrapper.Wrapper
	mu     sync.Mutex
	config Config
	cancel context.CancelFunc
	ctx    context.Context

	// evacuatedFleets tracks fleets we sent to evacuate, keyed by origin celestial ID.
	evacuatedFleets map[ogame.CelestialID]evacuatedFleet
}

type evacuatedFleet struct {
	FleetID     ogame.FleetID
	AttackEvent ogame.AttackEvent
	SentAt      time.Time
}

// New creates a new Defender with the given bot wrapper and config.
func New(bot wrapper.Wrapper, cfg Config) *Defender {
	cfg.Validate()
	return &Defender{
		bot:             bot,
		config:          cfg,
		evacuatedFleets: make(map[ogame.CelestialID]evacuatedFleet),
	}
}

// GetConfig returns a copy of the current configuration.
func (d *Defender) GetConfig() Config {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.config
}

// SetConfig updates the configuration. If the defender is running, it restarts with the new config.
func (d *Defender) SetConfig(cfg Config) {
	cfg.Validate()
	d.mu.Lock()
	wasRunning := d.cancel != nil
	d.mu.Unlock()

	if wasRunning {
		d.Stop()
	}

	d.mu.Lock()
	d.config = cfg
	d.mu.Unlock()

	if wasRunning && cfg.Enabled {
		d.Start()
	}
}

// Start begins the defender monitoring loop.
func (d *Defender) Start() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cancel != nil {
		return // already running
	}

	if !d.config.Enabled {
		log.Println("[Defender] Not enabled, skipping start")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.ctx = ctx
	d.cancel = cancel

	log.Printf("[Defender] Starting with check interval %ds, evac speed %.0f%%, evac %ds before attack",
		d.config.CheckIntervalSec, d.config.EvacSpeed.Float64()*10, d.config.EvacSecondsBefore)

	go d.loop(ctx)
}

// Stop stops the defender monitoring loop.
func (d *Defender) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cancel != nil {
		log.Println("[Defender] Stopping")
		d.cancel()
		d.cancel = nil
		d.ctx = nil
	}
}

// IsRunning returns whether the defender is currently active.
func (d *Defender) IsRunning() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cancel != nil
}

func (d *Defender) loop(ctx context.Context) {
	d.mu.Lock()
	interval := time.Duration(d.config.CheckIntervalSec) * time.Second
	d.mu.Unlock()

	// Perform an initial check immediately.
	d.checkAndDefend()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Defender] Loop stopped")
			return
		case <-ticker.C:
			d.checkAndDefend()
		}
	}
}

func (d *Defender) checkAndDefend() {
	d.mu.Lock()
	cfg := d.config
	d.mu.Unlock()

	if !cfg.Enabled {
		return
	}

	log.Println("[Defender] Checking for attacks...")

	attacks, err := d.bot.GetAttacks()
	if err != nil {
		log.Printf("[Defender] Error getting attacks: %v", err)
		return
	}

	if len(attacks) == 0 {
		log.Println("[Defender] No attacks detected")
		// Check if any evacuated fleets can be recalled.
		d.checkRecalls(cfg)
		return
	}

	log.Printf("[Defender] %d attack(s) detected!", len(attacks))

	for _, attack := range attacks {
		d.handleAttack(cfg, attack)
	}
}

func (d *Defender) handleAttack(cfg Config, attack ogame.AttackEvent) {
	// Only handle attacks and destroy missions targeted at our celestials.
	if attack.MissionType != ogame.Attack && attack.MissionType != ogame.Destroy &&
		attack.MissionType != ogame.GroupedAttack && attack.MissionType != ogame.MissileAttack {
		return
	}

	// Find the celestial being attacked.
	destCoord := attack.Destination
	celestial, err := d.bot.GetCachedCelestial(destCoord)
	if err != nil {
		log.Printf("[Defender] Could not find celestial at %s: %v", destCoord.String(), err)
		return
	}

	celestialID := celestial.GetID()

	// Check if we already evacuated from this celestial.
	d.mu.Lock()
	if _, exists := d.evacuatedFleets[celestialID]; exists {
		d.mu.Unlock()
		log.Printf("[Defender] Already evacuated from %s, skipping", destCoord.String())
		return
	}
	d.mu.Unlock()

	// Calculate time until attack.
	timeUntilAttack := time.Until(attack.ArrivalTime)
	evacBefore := time.Duration(cfg.EvacSecondsBefore) * time.Second

	if timeUntilAttack <= 0 {
		log.Printf("[Defender] Attack on %s already arrived, too late to evacuate", destCoord.String())
		return
	}

	log.Printf("[Defender] Attack on %s by %s (ID:%d), arrives in %s (mission: %s)",
		destCoord.String(), attack.AttackerName, attack.AttackerID,
		timeUntilAttack.Round(time.Second), attack.MissionType.String())

	// If the attack is still far away, schedule evacuation for later.
	waitTime := timeUntilAttack - evacBefore
	if waitTime > 0 {
		log.Printf("[Defender] Scheduling evacuation of %s in %s", destCoord.String(), waitTime.Round(time.Second))
		go func() {
			d.mu.Lock()
			ctx := d.ctx
			d.mu.Unlock()
			if ctx == nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(waitTime):
				d.evacuate(cfg, celestial, attack)
			}
		}()
		return
	}

	// Attack is imminent, evacuate now.
	d.evacuate(cfg, celestial, attack)
}

func (d *Defender) evacuate(cfg Config, celestial wrapper.Celestial, attack ogame.AttackEvent) {
	celestialID := celestial.GetID()
	destCoord := attack.Destination

	// Double-check we haven't already evacuated.
	d.mu.Lock()
	if _, exists := d.evacuatedFleets[celestialID]; exists {
		d.mu.Unlock()
		return
	}
	d.mu.Unlock()

	log.Printf("[Defender] Evacuating fleet from %s NOW!", destCoord.String())

	// Find a safe destination (a different celestial we own).
	safeDest, err := d.findSafeDestination(celestial)
	if err != nil {
		log.Printf("[Defender] Cannot find safe destination for %s: %v", destCoord.String(), err)
		return
	}

	log.Printf("[Defender] Safe destination: %s", safeDest.String())

	// Send all ships with all resources to the safe destination.
	fleet, err := wrapper.NewFleetBuilder(d.bot).
		SetOrigin(celestial).
		SetDestination(safeDest).
		SetAllShips().
		SetAllResources().
		SetSpeed(cfg.EvacSpeed).
		SetMission(ogame.Park).
		SendNow()

	if err != nil {
		log.Printf("[Defender] Failed to evacuate from %s: %v", destCoord.String(), err)
		return
	}

	log.Printf("[Defender] Evacuation fleet %d sent from %s to %s", fleet.ID, destCoord.String(), safeDest.String())

	d.mu.Lock()
	d.evacuatedFleets[celestialID] = evacuatedFleet{
		FleetID:     fleet.ID,
		AttackEvent: attack,
		SentAt:      time.Now(),
	}
	d.mu.Unlock()

	// Schedule recall if configured.
	if cfg.RecallSecondsAfter > 0 {
		recallDelay := time.Until(attack.ArrivalTime) + time.Duration(cfg.RecallSecondsAfter)*time.Second
		if recallDelay < 0 {
			recallDelay = time.Duration(cfg.RecallSecondsAfter) * time.Second
		}
		go func() {
			d.mu.Lock()
			ctx := d.ctx
			d.mu.Unlock()
			if ctx == nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(recallDelay):
				d.tryRecall(cfg, celestialID, fleet.ID)
			}
		}()
	}
}

func (d *Defender) findSafeDestination(origin wrapper.Celestial) (ogame.Coordinate, error) {
	celestials := d.bot.GetCachedCelestials()
	originCoord := origin.GetCoordinate()

	// Prefer a moon if the origin is a planet (and vice versa) at the same position.
	for _, c := range celestials {
		coord := c.GetCoordinate()
		if coord.Galaxy == originCoord.Galaxy &&
			coord.System == originCoord.System &&
			coord.Position == originCoord.Position &&
			coord.Type != originCoord.Type {
			return coord, nil
		}
	}

	// Otherwise pick the closest different celestial.
	var bestDest ogame.Coordinate
	bestDist := int64(-1)
	for _, c := range celestials {
		if c.GetID() == origin.GetID() {
			continue
		}
		dist := d.bot.Distance(originCoord, c.GetCoordinate())
		if bestDist == -1 || dist < bestDist {
			bestDist = dist
			bestDest = c.GetCoordinate()
		}
	}

	if bestDist == -1 {
		return ogame.Coordinate{}, fmt.Errorf("no other celestial available")
	}

	return bestDest, nil
}

func (d *Defender) checkRecalls(cfg Config) {
	if cfg.RecallSecondsAfter <= 0 {
		return
	}

	d.mu.Lock()
	fleetsCopy := make(map[ogame.CelestialID]evacuatedFleet, len(d.evacuatedFleets))
	for k, v := range d.evacuatedFleets {
		fleetsCopy[k] = v
	}
	d.mu.Unlock()

	for celestialID, evac := range fleetsCopy {
		// Check if enough time has passed since the attack was supposed to arrive.
		timeSinceAttack := time.Since(evac.AttackEvent.ArrivalTime)
		if timeSinceAttack < time.Duration(cfg.RecallSecondsAfter)*time.Second {
			continue
		}

		d.tryRecall(cfg, celestialID, evac.FleetID)
	}
}

func (d *Defender) tryRecall(cfg Config, celestialID ogame.CelestialID, fleetID ogame.FleetID) {
	// Check if there are still attacks on this celestial.
	attacks, err := d.bot.GetAttacks()
	if err != nil {
		log.Printf("[Defender] Error checking attacks for recall: %v", err)
		return
	}

	celestial, err := d.bot.GetCachedCelestial(celestialID)
	if err != nil {
		log.Printf("[Defender] Error getting celestial %d for recall: %v", celestialID, err)
		return
	}

	coord := celestial.GetCoordinate()

	for _, attack := range attacks {
		if attack.Destination.Equal(coord) {
			log.Printf("[Defender] Still under attack at %s, not recalling fleet %d", coord.String(), fleetID)
			return
		}
	}

	// No more attacks, recall the fleet.
	log.Printf("[Defender] No more attacks on %s, recalling fleet %d", coord.String(), fleetID)
	err = d.bot.CancelFleet(fleetID)
	if err != nil {
		log.Printf("[Defender] Failed to recall fleet %d: %v", fleetID, err)
	} else {
		log.Printf("[Defender] Fleet %d recalled successfully", fleetID)
	}

	d.mu.Lock()
	delete(d.evacuatedFleets, celestialID)
	d.mu.Unlock()
}

// selectCheckCelestial returns the celestial to use for attack checking based on config.
func (d *Defender) selectCheckCelestial(cfg Config) (wrapper.Celestial, error) {
	if cfg.Origin.Random {
		celestials := d.bot.GetCachedCelestials()
		if len(celestials) == 0 {
			return nil, fmt.Errorf("no celestials available")
		}
		return celestials[rand.Intn(len(celestials))], nil
	}
	return d.bot.GetCachedCelestial(cfg.Origin.CelestialID)
}
