package ogame

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

const (
	SMALL_CARGO = iota
	LARGE_CARGO
	LIGHT_FIGHTER
	HEAVY_FIGHTER
	CRUISER
	BATTLESHIP
	COLONY_SHIP
	RECYCLER
	ESPIONAGE_PROBE
	BOMBER
	SOLAR_SATELLITE
	DESTROYER
	DEATHSTAR
	BATTLECRUISER
	ROCKET_LAUNCHER
	LIGHT_LASER
	HEAVY_LASER
	GAUSS_CANNON
	ION_CANNON
	PLASMA_TURRET
	SMALL_SHIELD_DOME
	LARGE_SHIELD_DOME
)

func isAlive(unit *CombatUnit) bool {
	return getUnitHull(unit) > 0
}

type CombatUnit struct {
	PackedInfos uint64
}

const (
	MAX_ARMOUR_LEVEL uint64 = 36
	MAX_SHIELD_LEVEL uint64 = 42
	ID_MASK          uint64 = 31
	SHIELD_MASK      uint64 = 8388576
	HULL_MASK        uint64 = 35184363700224
)

func getUnitId(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & ID_MASK) >> 0
}

func getUnitShield(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & SHIELD_MASK) >> 5
}

func getUnitHull(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & HULL_MASK) >> 23
}

func setUnitId(unit *CombatUnit, id uint64) {
	unit.PackedInfos &= ^ID_MASK
	unit.PackedInfos |= id << 0
}

func setUnitShield(unit *CombatUnit, shield uint64) {
	unit.PackedInfos &= ^SHIELD_MASK
	unit.PackedInfos |= shield << 5
}

func setUnitHull(unit *CombatUnit, hull uint64) {
	unit.PackedInfos &= ^HULL_MASK
	unit.PackedInfos |= hull << 23
}

type price struct {
	Metal     int
	Crystal   int
	Deuterium int
}

func (p *price) Add(n price) {
	p.Metal += n.Metal
	p.Crystal += n.Crystal
	p.Deuterium += n.Deuterium
}

func getUnitPrice(unitId uint64) price {
	switch unitId {
	case SMALL_CARGO:
		return price{2000, 2000, 0}
	case LARGE_CARGO:
		return price{6000, 6000, 0}
	case LIGHT_FIGHTER:
		return price{3000, 1000, 0}
	case HEAVY_FIGHTER:
		return price{6000, 4000, 0}
	case CRUISER:
		return price{20000, 7000, 2000}
	case BATTLESHIP:
		return price{45000, 15000, 0}
	case COLONY_SHIP:
		return price{10000, 20000, 10000}
	case RECYCLER:
		return price{10000, 6000, 2000}
	case ESPIONAGE_PROBE:
		return price{0, 1000, 0}
	case BOMBER:
		return price{50000, 25000, 15000}
	case SOLAR_SATELLITE:
		return price{0, 2000, 500}
	case DESTROYER:
		return price{60000, 50000, 15000}
	case DEATHSTAR:
		return price{5000000, 4000000, 1000000}
	case BATTLECRUISER:
		return price{30000, 40000, 15000}
	case ROCKET_LAUNCHER:
		return price{2000, 0, 0}
	case LIGHT_LASER:
		return price{1500, 500, 0}
	case HEAVY_LASER:
		return price{6000, 2000, 0}
	case GAUSS_CANNON:
		return price{20000, 15000, 2000}
	case ION_CANNON:
		return price{2000, 6000, 0}
	case PLASMA_TURRET:
		return price{50000, 50000, 30000}
	case SMALL_SHIELD_DOME:
		return price{10000, 10000, 0}
	case LARGE_SHIELD_DOME:
		return price{50000, 50000, 0}
	}
	return price{0, 0, 0}
}

func getUnitBaseShield(unitId uint64) int {
	switch unitId {
	case SMALL_CARGO:
		return 10
	case LARGE_CARGO:
		return 25
	case LIGHT_FIGHTER:
		return 10
	case HEAVY_FIGHTER:
		return 25
	case CRUISER:
		return 50
	case BATTLESHIP:
		return 200
	case COLONY_SHIP:
		return 100
	case RECYCLER:
		return 10
	case ESPIONAGE_PROBE:
		return 1 // 0.01
	case BOMBER:
		return 500
	case SOLAR_SATELLITE:
		return 1
	case DESTROYER:
		return 500
	case DEATHSTAR:
		return 50000
	case BATTLECRUISER:
		return 400
	case ROCKET_LAUNCHER:
		return 20
	case LIGHT_LASER:
		return 25
	case HEAVY_LASER:
		return 100
	case GAUSS_CANNON:
		return 200
	case ION_CANNON:
		return 500
	case PLASMA_TURRET:
		return 300
	case SMALL_SHIELD_DOME:
		return 2000
	case LARGE_SHIELD_DOME:
		return 10000
	}
	return 0
}

func getUnitBaseWeapon(unitId uint64) uint64 {
	switch unitId {
	case SMALL_CARGO:
		return 5
	case LARGE_CARGO:
		return 5
	case LIGHT_FIGHTER:
		return 50
	case HEAVY_FIGHTER:
		return 150
	case CRUISER:
		return 400
	case BATTLESHIP:
		return 1000
	case COLONY_SHIP:
		return 50
	case RECYCLER:
		return 1
	case ESPIONAGE_PROBE:
		return 1 // 0.01
	case BOMBER:
		return 1000
	case SOLAR_SATELLITE:
		return 1
	case DESTROYER:
		return 2000
	case DEATHSTAR:
		return 200000
	case BATTLECRUISER:
		return 700
	case ROCKET_LAUNCHER:
		return 80
	case LIGHT_LASER:
		return 100
	case HEAVY_LASER:
		return 250
	case GAUSS_CANNON:
		return 1100
	case ION_CANNON:
		return 150
	case PLASMA_TURRET:
		return 3000
	case SMALL_SHIELD_DOME:
		return 1
	case LARGE_SHIELD_DOME:
		return 1
	}
	return 0
}

func getRapidFireAgainst(unit *CombatUnit, targetUnit *CombatUnit) int {
	rf := 0
	switch getUnitId(unit) {
	case SMALL_CARGO:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case LARGE_CARGO:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case LIGHT_FIGHTER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case HEAVY_FIGHTER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		case SMALL_CARGO:
			rf = 3
			break
		}
		break
	case CRUISER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		case LIGHT_FIGHTER:
			rf = 6
			break
		case ROCKET_LAUNCHER:
			rf = 10
			break
		}
		break
	case BATTLESHIP:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case COLONY_SHIP:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case RECYCLER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		}
		break
	case BOMBER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		case ION_CANNON:
			rf = 10
			break
		case ROCKET_LAUNCHER:
			rf = 20
			break
		case LIGHT_LASER:
			rf = 20
			break
		case HEAVY_LASER:
			rf = 10
			break
		}
		break
	case DESTROYER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		case LIGHT_LASER:
			rf = 10
			break
		case BATTLECRUISER:
			rf = 2
			break
		}
		break
	case DEATHSTAR:
		switch getUnitId(targetUnit) {
		case SMALL_CARGO:
			rf = 250
			break
		case LARGE_CARGO:
			rf = 250
			break
		case LIGHT_FIGHTER:
			rf = 200
			break
		case HEAVY_FIGHTER:
			rf = 100
			break
		case CRUISER:
			rf = 33
			break
		case BATTLESHIP:
			rf = 30
			break
		case COLONY_SHIP:
			rf = 250
			break
		case RECYCLER:
			rf = 250
			break
		case ESPIONAGE_PROBE:
			rf = 1250
			break
		case SOLAR_SATELLITE:
			rf = 1250
			break
		case BOMBER:
			rf = 25
			break
		case DESTROYER:
			rf = 5
			break
		case ROCKET_LAUNCHER:
			rf = 200
			break
		case LIGHT_LASER:
			rf = 200
			break
		case HEAVY_LASER:
			rf = 100
			break
		case GAUSS_CANNON:
			rf = 50
			break
		case ION_CANNON:
			rf = 100
			break
		case BATTLECRUISER:
			rf = 15
			break
		}
		break
	case BATTLECRUISER:
		switch getUnitId(targetUnit) {
		case ESPIONAGE_PROBE:
			rf = 5
			break
		case SOLAR_SATELLITE:
			rf = 5
			break
		case SMALL_CARGO:
			rf = 3
			break
		case LARGE_CARGO:
			rf = 3
			break
		case HEAVY_FIGHTER:
			rf = 4
			break
		case CRUISER:
			rf = 4
			break
		case BATTLESHIP:
			rf = 7
			break
		}
		break
	}
	return rf
}

func getUnitName(unitId uint64) string {
	switch unitId {
	case SMALL_CARGO:
		return "Small cargo"
	case LARGE_CARGO:
		return "Large cargo"
	case LIGHT_FIGHTER:
		return "Light fighter"
	case HEAVY_FIGHTER:
		return "Heavy fighter"
	case CRUISER:
		return "Cruiser"
	case BATTLESHIP:
		return "Battleship"
	case COLONY_SHIP:
		return "Colony ship"
	case RECYCLER:
		return "Recycler"
	case ESPIONAGE_PROBE:
		return "Expionage probe"
	case BOMBER:
		return "Bomber"
	case SOLAR_SATELLITE:
		return "Solar satellite"
	case DESTROYER:
		return "Destroyer"
	case DEATHSTAR:
		return "Deathstar"
	case BATTLECRUISER:
		return "Battlecruiser"
	case ROCKET_LAUNCHER:
		return "Rocket launcher"
	case LIGHT_LASER:
		return "Light laser"
	case HEAVY_LASER:
		return "Heavy laser"
	case GAUSS_CANNON:
		return "Gauss cannon"
	case ION_CANNON:
		return "Ion cannon"
	case PLASMA_TURRET:
		return "Plasma turret"
	case SMALL_SHIELD_DOME:
		return "Small shield dome"
	case LARGE_SHIELD_DOME:
		return "Large shield dome"
	}
	return ""
}

func getUnitWeaponPower(unitId uint64, weaponTechno int) uint64 {
	return uint64(float64(getUnitBaseWeapon(unitId)) * (1 + 0.1*float64(weaponTechno)))
}

func getUnitInitialShield(unitId uint64, shieldTechno int) uint64 {
	return uint64(float64(getUnitBaseShield(unitId)) * (1 + 0.1*float64(shieldTechno)))
}

func getUnitInitialHullPlating(armourTechno, metalPrice, crystalPrice int) uint64 {
	return uint64((1 + (float64(armourTechno) / 10)) * (float64(metalPrice+crystalPrice) / 10))
}

func newUnit(entity *entity, unitId uint64) CombatUnit {
	var unit CombatUnit
	setUnitId(&unit, unitId)
	unitPrice := getUnitPrice(unitId)
	setUnitHull(&unit, getUnitInitialHullPlating(entity.Armour, unitPrice.Metal, unitPrice.Crystal))
	setUnitShield(&unit, getUnitInitialShield(unitId, entity.Shield))
	return unit
}

type entity struct {
	Weapon          int
	Shield          int
	Armour          int
	Combustion      int
	Impulse         int
	Hyperspace      int
	SmallCargo      int
	LargeCargo      int
	LightFighter    int
	HeavyFighter    int
	Cruiser         int
	Battleship      int
	ColonyShip      int
	Recycler        int
	EspionageProbe  int
	Bomber          int
	SolarSatellite  int
	Destroyer       int
	Deathstar       int
	Battlecruiser   int
	RocketLauncher  int
	LightLaser      int
	HeavyLaser      int
	GaussCannon     int
	IonCannon       int
	PlasmaTurret    int
	SmallShieldDome int
	LargeShieldDome int
	TotalUnits      int
	Units           []CombatUnit
	Losses          price
}

func (e *entity) Init() {
	e.Reset()
	idx := 0
	for i := 0; i < e.SmallCargo; i++ {
		e.Units[idx] = newUnit(e, SMALL_CARGO)
		idx++
	}
	for i := 0; i < e.LargeCargo; i++ {
		e.Units[idx] = newUnit(e, LARGE_CARGO)
		idx++
	}
	for i := 0; i < e.LightFighter; i++ {
		e.Units[idx] = newUnit(e, LIGHT_FIGHTER)
		idx++
	}
	for i := 0; i < e.HeavyFighter; i++ {
		e.Units[idx] = newUnit(e, HEAVY_FIGHTER)
		idx++
	}
	for i := 0; i < e.Cruiser; i++ {
		e.Units[idx] = newUnit(e, CRUISER)
		idx++
	}
	for i := 0; i < e.Battleship; i++ {
		e.Units[idx] = newUnit(e, BATTLESHIP)
		idx++
	}
	for i := 0; i < e.ColonyShip; i++ {
		e.Units[idx] = newUnit(e, COLONY_SHIP)
		idx++
	}
	for i := 0; i < e.Recycler; i++ {
		e.Units[idx] = newUnit(e, RECYCLER)
		idx++
	}
	for i := 0; i < e.EspionageProbe; i++ {
		e.Units[idx] = newUnit(e, ESPIONAGE_PROBE)
		idx++
	}
	for i := 0; i < e.Bomber; i++ {
		e.Units[idx] = newUnit(e, BOMBER)
		idx++
	}
	for i := 0; i < e.SolarSatellite; i++ {
		e.Units[idx] = newUnit(e, SOLAR_SATELLITE)
		idx++
	}
	for i := 0; i < e.Destroyer; i++ {
		e.Units[idx] = newUnit(e, DESTROYER)
		idx++
	}
	for i := 0; i < e.Deathstar; i++ {
		e.Units[idx] = newUnit(e, DEATHSTAR)
		idx++
	}
	for i := 0; i < e.Battlecruiser; i++ {
		e.Units[idx] = newUnit(e, BATTLECRUISER)
		idx++
	}
	for i := 0; i < e.RocketLauncher; i++ {
		e.Units[idx] = newUnit(e, ROCKET_LAUNCHER)
		idx++
	}
	for i := 0; i < e.LightLaser; i++ {
		e.Units[idx] = newUnit(e, LIGHT_LASER)
		idx++
	}
	for i := 0; i < e.HeavyLaser; i++ {
		e.Units[idx] = newUnit(e, HEAVY_LASER)
		idx++
	}
	for i := 0; i < e.GaussCannon; i++ {
		e.Units[idx] = newUnit(e, GAUSS_CANNON)
		idx++
	}
	for i := 0; i < e.IonCannon; i++ {
		e.Units[idx] = newUnit(e, ION_CANNON)
		idx++
	}
	for i := 0; i < e.PlasmaTurret; i++ {
		e.Units[idx] = newUnit(e, PLASMA_TURRET)
		idx++
	}
	for i := 0; i < e.SmallShieldDome; i++ {
		e.Units[idx] = newUnit(e, SMALL_SHIELD_DOME)
		idx++
	}
	for i := 0; i < e.LargeShieldDome; i++ {
		e.Units[idx] = newUnit(e, LARGE_SHIELD_DOME)
		idx++
	}
}

func newEntity() *entity {
	return new(entity)
}

type combatSimulator struct {
	Attacker      entity
	Defender      entity
	MaxRounds     int
	Rounds        int
	FleetToDebris float64
	Winner        string
	IsLogging     bool
	Logs          string
	Debris        price
}

func (simulator *combatSimulator) hasExploded(entity *entity, defendingUnit *CombatUnit) bool {
	exploded := false
	unitPrice := getUnitPrice(getUnitId(defendingUnit))
	hullPercentage := float64(getUnitHull(defendingUnit)) / float64(getUnitInitialHullPlating(entity.Armour, unitPrice.Metal, unitPrice.Crystal))
	if hullPercentage <= 0.7 {
		probabilityOfExploding := 1.0 - hullPercentage
		dice := rand.Float64()
		msg := ""
		if simulator.IsLogging {
			msg += fmt.Sprintf("probability of exploding of %1.3f%%: dice value of %1.3f comparing with %1.3f: ", probabilityOfExploding*100, dice, 1-probabilityOfExploding)
		}
		if dice >= 1-probabilityOfExploding {
			exploded = true
			if simulator.IsLogging {
				msg += "unit exploded."
			}
		} else {
			if simulator.IsLogging {
				msg += "unit didn't explode."
			}
		}
		if simulator.IsLogging {
			simulator.Logs += msg + "\n"
		}
	}
	return exploded
}

func (simulator *combatSimulator) getAnotherShot(unit, targetUnit *CombatUnit) bool {
	rapidFire := true
	rf := getRapidFireAgainst(unit, targetUnit)
	msg := ""
	if rf > 0 {
		chance := float64(rf-1) / float64(rf)
		dice := rand.Float64()
		if simulator.IsLogging {
			msg += fmt.Sprintf("dice was %1.3f, comparing with %1.3f: ", dice, chance)
		}
		if dice <= chance {
			if simulator.IsLogging {
				msg += fmt.Sprintf("%s gets another shot.", getUnitName(getUnitId(unit)))
			}
		} else {
			if simulator.IsLogging {
				msg += fmt.Sprintf("%s does not get another shot.", getUnitName(getUnitId(unit)))
			}
			rapidFire = false
		}
	} else {
		if simulator.IsLogging {
			msg += fmt.Sprintf("%s doesn't have rapid fire against %s.", getUnitName(getUnitId(unit)), getUnitName(getUnitId(targetUnit)))
		}
		rapidFire = false
	}
	if simulator.IsLogging {
		simulator.Logs += msg + "\n"
	}
	return rapidFire
}

func (simulator *combatSimulator) attack(attacker *entity, attackingUnit *CombatUnit, defender *entity, defendingUnit *CombatUnit) {
	if simulator.IsLogging {
		simulator.Logs += fmt.Sprintf("%s fires at %s; ", getUnitName(getUnitId(attackingUnit)), getUnitName(getUnitId(defendingUnit)))
	}

	weapon := getUnitWeaponPower(getUnitId(attackingUnit), attacker.Weapon)
	// Check for shot bounce
	if float64(weapon) < 0.01*float64(getUnitShield(defendingUnit)) {
		if simulator.IsLogging {
			simulator.Logs += "shot bounced\n"
		}
		return
	}

	// Attack target
	currentHull := getUnitHull(defendingUnit)
	currentShield := getUnitShield(defendingUnit)
	if currentShield < weapon {
		weapon -= currentShield
		setUnitShield(defendingUnit, 0)
		if (int64)(currentHull-weapon) < 0 {
			setUnitHull(defendingUnit, 0)
		} else {
			setUnitHull(defendingUnit, currentHull-weapon)
		}
	} else {
		setUnitShield(defendingUnit, currentShield-weapon)
	}
	if simulator.IsLogging {
		simulator.Logs += fmt.Sprintf("result is %s %d %d\n", getUnitName(getUnitId(defendingUnit)), getUnitHull(defendingUnit), getUnitShield(defendingUnit))
	}

	// Check for explosion
	if isAlive(defendingUnit) {
		if simulator.hasExploded(defender, defendingUnit) {
			setUnitHull(defendingUnit, 0)
		}
	}
}

func (simulator *combatSimulator) unitsFires(attacker, defender *entity) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < attacker.TotalUnits; i++ {
		unit := attacker.Units[i]
		rapidFire := true
		for rapidFire {
			if defender.TotalUnits == 0 {
				break
			}
			targetUnit := &defender.Units[rand.Intn(defender.TotalUnits)]
			rapidFire = simulator.getAnotherShot(&unit, targetUnit)
			if isAlive(targetUnit) {
				simulator.attack(attacker, &unit, defender, targetUnit)
			}
		}
	}
}

func (simulator *combatSimulator) attackerFires() {
	if simulator.Defender.TotalUnits <= 0 {
		return
	}
	simulator.unitsFires(&simulator.Attacker, &simulator.Defender)
}

func (simulator *combatSimulator) defenderFires() {
	simulator.unitsFires(&simulator.Defender, &simulator.Attacker)
}

func isShip(unit *CombatUnit) bool {
	switch getUnitId(unit) {
	case SMALL_CARGO:
		return true
	case LARGE_CARGO:
		return true
	case LIGHT_FIGHTER:
		return true
	case HEAVY_FIGHTER:
		return true
	case CRUISER:
		return true
	case BATTLESHIP:
		return true
	case COLONY_SHIP:
		return true
	case RECYCLER:
		return true
	case ESPIONAGE_PROBE:
		return true
	case BOMBER:
		return true
	case SOLAR_SATELLITE:
		return true
	case DESTROYER:
		return true
	case DEATHSTAR:
		return true
	case BATTLECRUISER:
		return true
	}
	return false
}

func (simulator *combatSimulator) RemoveDestroyedUnits() {
	l := simulator.Defender.TotalUnits
	for i := l - 1; i >= 0; i-- {
		unit := &simulator.Defender.Units[i]
		if getUnitHull(unit) <= 0 {
			unitPrice := getUnitPrice(getUnitId(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Defender.Losses.Add(unitPrice)
			simulator.Defender.Units[i] = simulator.Defender.Units[simulator.Defender.TotalUnits-1]
			simulator.Defender.TotalUnits--
			//simulator.Defender.Units = simulator.Defender.Units[:len(simulator.Defender.Units)-1]
			if simulator.IsLogging {
				simulator.Logs += fmt.Sprintf("%s lost all its integrity, remove from battle\n", getUnitName(getUnitId(unit)))
			}
		}
	}
	l = simulator.Attacker.TotalUnits
	for i := l - 1; i >= 0; i-- {
		unit := &simulator.Attacker.Units[i]
		if getUnitHull(unit) <= 0 {
			unitPrice := getUnitPrice(getUnitId(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Attacker.Losses.Add(unitPrice)
			simulator.Attacker.Units[i] = simulator.Attacker.Units[simulator.Attacker.TotalUnits-1]
			simulator.Attacker.TotalUnits--
			//simulator.Attacker.Units = simulator.Attacker.Units[:len(simulator.Attacker.Units)-1]
			if simulator.IsLogging {
				simulator.Logs += fmt.Sprintf("%s lost all its integrity, remove from battle\n", getUnitName(getUnitId(unit)))
			}
		}
	}
}

func (simulator *combatSimulator) RestoreShields() {
	for i := 0; i < simulator.Attacker.TotalUnits; i++ {
		unit := &simulator.Attacker.Units[i]
		setUnitShield(unit, getUnitInitialShield(getUnitId(unit), simulator.Attacker.Shield))
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("%s still has integrity, restore its shield\n", getUnitName(getUnitId(unit)))
		}
	}
	for i := 0; i < simulator.Defender.TotalUnits; i++ {
		unit := &simulator.Defender.Units[i]
		setUnitShield(unit, getUnitInitialShield(getUnitId(unit), simulator.Defender.Shield))
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("%s still has integrity, restore its shield\n", getUnitName(getUnitId(unit)))
		}
	}
}

func (simulator *combatSimulator) IsCombatDone() bool {
	return simulator.Attacker.TotalUnits <= 0 || simulator.Defender.TotalUnits <= 0
}

func (simulator *combatSimulator) GetMoonchance() int {
	debris := float64(simulator.Debris.Metal) + float64(simulator.Debris.Crystal)
	return int(math.Min(debris/100000.0, 20.0))
}

func (simulator *combatSimulator) PrintWinner() {
	if simulator.Defender.TotalUnits <= 0 && simulator.Attacker.TotalUnits <= 0 {
		simulator.Winner = "draw"
		if simulator.IsLogging {
			simulator.Logs += "The battle ended draw.\n"
		}
	} else if simulator.Attacker.TotalUnits <= 0 {
		simulator.Winner = "defender"
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("The battle ended after %d rounds with %s winning\n", simulator.Rounds, simulator.Winner)
		}
	} else if simulator.Defender.TotalUnits <= 0 {
		simulator.Winner = "attacker"
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("The battle ended after %d rounds with %s winning\n", simulator.Rounds, simulator.Winner)
		}
	} else {
		simulator.Winner = "draw"
		if simulator.IsLogging {
			simulator.Logs += "The battle ended draw.\n"
		}
	}
}

func (simulator *combatSimulator) Simulate() {
	simulator.Attacker.Init()
	simulator.Defender.Init()
	for currentRound := 1; currentRound <= simulator.MaxRounds; currentRound++ {
		simulator.Rounds = currentRound
		if simulator.IsLogging {
			simulator.Logs += strings.Repeat("-", 80) + "\n"
			simulator.Logs += "ROUND " + strconv.Itoa(currentRound) + "\n"
			simulator.Logs += strings.Repeat("-", 80) + "\n"
		}
		simulator.attackerFires()
		simulator.defenderFires()
		simulator.RemoveDestroyedUnits()
		simulator.RestoreShields()
		if simulator.IsCombatDone() {
			break
		}
	}
	simulator.PrintWinner()
}

func newCombatSimulator(attacker *entity, defender *entity) *combatSimulator {
	cs := new(combatSimulator)
	cs.Attacker = *attacker
	cs.Defender = *defender
	cs.IsLogging = false
	cs.MaxRounds = 6
	cs.FleetToDebris = 0.3
	return cs
}

type Config struct {
	IsLogging   bool
	Simulations int
	Workers     int
	Attacker    attackerInfo
	Defender    defenderInfo
}

type attackerInfo struct {
	Weapon         int
	Shield         int
	Armour         int
	SmallCargo     int
	LargeCargo     int
	LightFighter   int
	HeavyFighter   int
	Cruiser        int
	Battleship     int
	ColonyShip     int
	Recycler       int
	EspionageProbe int
	Bomber         int
	Destroyer      int
	Deathstar      int
	Battlecruiser  int
}

type defenderInfo struct {
	Weapon          int
	Shield          int
	Armour          int
	SmallCargo      int
	LargeCargo      int
	LightFighter    int
	HeavyFighter    int
	Cruiser         int
	Battleship      int
	ColonyShip      int
	Recycler        int
	EspionageProbe  int
	Bomber          int
	SolarSatellite  int
	Destroyer       int
	Deathstar       int
	Battlecruiser   int
	RocketLauncher  int
	LightLaser      int
	HeavyLaser      int
	GaussCannon     int
	IonCannon       int
	PlasmaTurret    int
	SmallShieldDome int
	LargeShieldDome int
}

func printResult(result SimulatorResult) {
	fmt.Println(fmt.Sprintf("| Results (%d simulations | ~%d rounds)", result.Simulations, result.Rounds))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	data2 := [][]string{
		{"Attackers win", fmt.Sprintf("%d%%", result.AttackerWin)},
		{"Defenders win", fmt.Sprintf("%d%%", result.DefenderWin)},
		{"Draw", fmt.Sprintf("%d%%", result.Draw)},
	}
	table.AppendBulk(data2)
	table.Render()

	fmt.Println("")
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Metal", "Crystal", "Deuterium", "Recycler", "Moonchance"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	data1 := [][]string{
		{
			"Attacker losses",
			fmt.Sprintf("%d", result.AttackerLosses.Metal),
			fmt.Sprintf("%d", result.AttackerLosses.Crystal),
			fmt.Sprintf("%d", result.AttackerLosses.Deuterium),
			"", ""},
		{
			"Defender losses",
			fmt.Sprintf("%d", result.DefenderLosses.Metal),
			fmt.Sprintf("%d", result.DefenderLosses.Crystal),
			fmt.Sprintf("%d", result.DefenderLosses.Deuterium),
			"", ""},
		{
			"Debris",
			fmt.Sprintf("%d", result.Debris.Metal),
			fmt.Sprintf("%d", result.Debris.Crystal),
			"",
			fmt.Sprintf("%d", result.Recycler),
			fmt.Sprintf("%d", result.Moonchance)},
	}
	table.AppendBulk(data1)
	table.Render()
}

func (e *entity) Reset() {
	e.Losses = price{Metal: 0, Crystal: 0, Deuterium: 0}
	e.TotalUnits = 0
	e.TotalUnits += e.SmallCargo
	e.TotalUnits += e.SmallCargo
	e.TotalUnits += e.LargeCargo
	e.TotalUnits += e.LightFighter
	e.TotalUnits += e.HeavyFighter
	e.TotalUnits += e.Cruiser
	e.TotalUnits += e.Battleship
	e.TotalUnits += e.ColonyShip
	e.TotalUnits += e.Recycler
	e.TotalUnits += e.EspionageProbe
	e.TotalUnits += e.Bomber
	e.TotalUnits += e.SolarSatellite
	e.TotalUnits += e.Destroyer
	e.TotalUnits += e.Deathstar
	e.TotalUnits += e.Battlecruiser
	e.TotalUnits += e.RocketLauncher
	e.TotalUnits += e.LightLaser
	e.TotalUnits += e.HeavyLaser
	e.TotalUnits += e.GaussCannon
	e.TotalUnits += e.IonCannon
	e.TotalUnits += e.PlasmaTurret
	e.TotalUnits += e.SmallShieldDome
	e.TotalUnits += e.LargeShieldDome
}

func Simulate(attackerParam Attacker, defenderParam Defender, params SimulatorParams) SimulatorResult {
	nbSimulations := params.Simulations

	attackerWin := 0
	defenderWin := 0
	draw := 0
	attackerLosses := price{}
	defenderLosses := price{}
	debris := price{}
	rounds := 0
	moonchance := 0

	attacker := newEntity()
	attacker.Weapon = attackerParam.Weapon
	attacker.Shield = attackerParam.Shield
	attacker.Armour = attackerParam.Armour
	attacker.SmallCargo = attackerParam.SmallCargo
	attacker.LargeCargo = attackerParam.LargeCargo
	attacker.LightFighter = attackerParam.LightFighter
	attacker.HeavyFighter = attackerParam.HeavyFighter
	attacker.Cruiser = attackerParam.Cruiser
	attacker.Battleship = attackerParam.Battleship
	attacker.ColonyShip = attackerParam.ColonyShip
	attacker.Recycler = attackerParam.Recycler
	attacker.EspionageProbe = attackerParam.EspionageProbe
	attacker.Bomber = attackerParam.Bomber
	attacker.SolarSatellite = 0
	attacker.Destroyer = attackerParam.Destroyer
	attacker.Deathstar = attackerParam.Deathstar
	attacker.Battlecruiser = attackerParam.Battlecruiser
	attacker.RocketLauncher = 0
	attacker.LightLaser = 0
	attacker.HeavyLaser = 0
	attacker.GaussCannon = 0
	attacker.IonCannon = 0
	attacker.PlasmaTurret = 0
	attacker.SmallShieldDome = 0
	attacker.LargeShieldDome = 0
	attacker.Reset()
	attacker.Units = make([]CombatUnit, attacker.TotalUnits+1, attacker.TotalUnits+1)

	defender := newEntity()
	defender.Weapon = defenderParam.Weapon
	defender.Shield = defenderParam.Shield
	defender.Armour = defenderParam.Armour
	defender.SmallCargo = defenderParam.SmallCargo
	defender.LargeCargo = defenderParam.LargeCargo
	defender.LightFighter = defenderParam.LightFighter
	defender.HeavyFighter = defenderParam.HeavyFighter
	defender.Cruiser = defenderParam.Cruiser
	defender.Battleship = defenderParam.Battleship
	defender.ColonyShip = defenderParam.ColonyShip
	defender.Recycler = defenderParam.Recycler
	defender.EspionageProbe = defenderParam.EspionageProbe
	defender.Bomber = defenderParam.Bomber
	defender.SolarSatellite = defenderParam.SolarSatellite
	defender.Destroyer = defenderParam.Destroyer
	defender.Deathstar = defenderParam.Deathstar
	defender.Battlecruiser = defenderParam.Battlecruiser
	defender.RocketLauncher = defenderParam.RocketLauncher
	defender.LightLaser = defenderParam.LightLaser
	defender.HeavyLaser = defenderParam.HeavyLaser
	defender.GaussCannon = defenderParam.GaussCannon
	defender.IonCannon = defenderParam.IonCannon
	defender.PlasmaTurret = defenderParam.PlasmaTurret
	defender.SmallShieldDome = defenderParam.SmallShieldDome
	defender.LargeShieldDome = defenderParam.LargeShieldDome
	defender.Reset()
	defender.Units = make([]CombatUnit, defender.TotalUnits+1, defender.TotalUnits+1)

	cs := newCombatSimulator(attacker, defender)
	cs.IsLogging = true

	for i := 0; i < nbSimulations; i++ {
		cs.Rounds = 1
		cs.Debris = price{}
		cs.Simulate()

		if cs.Winner == "attacker" {
			attackerWin++
		} else if cs.Winner == "defender" {
			defenderWin++
		} else {
			draw++
		}
		attackerLosses.Add(cs.Attacker.Losses)
		defenderLosses.Add(cs.Defender.Losses)
		debris.Add(cs.Debris)
		rounds += cs.Rounds
		moonchance += cs.GetMoonchance()
	}

	result := SimulatorResult{}
	result.Simulations = nbSimulations
	result.AttackerWin = round(float64(attackerWin) / float64(nbSimulations) * 100)
	result.DefenderWin = round(float64(defenderWin) / float64(nbSimulations) * 100)
	result.Draw = round(float64(draw) / float64(nbSimulations) * 100)
	result.Rounds = round(float64(rounds) / float64(nbSimulations))
	result.AttackerLosses = price{}
	result.AttackerLosses.Metal = int(float64(attackerLosses.Metal) / float64(nbSimulations))
	result.AttackerLosses.Crystal = int(float64(attackerLosses.Crystal) / float64(nbSimulations))
	result.AttackerLosses.Deuterium = int(float64(attackerLosses.Deuterium) / float64(nbSimulations))
	result.DefenderLosses = price{}
	result.DefenderLosses.Metal = int(float64(defenderLosses.Metal) / float64(nbSimulations))
	result.DefenderLosses.Crystal = int(float64(defenderLosses.Crystal) / float64(nbSimulations))
	result.DefenderLosses.Deuterium = int(float64(defenderLosses.Deuterium) / float64(nbSimulations))
	result.Debris = price{}
	result.Debris.Metal = int(float64(debris.Metal) / float64(nbSimulations))
	result.Debris.Crystal = int(float64(debris.Crystal) / float64(nbSimulations))
	result.Recycler = int(math.Ceil((float64(debris.Metal+debris.Crystal) / float64(nbSimulations)) / 20000.0))
	result.Moonchance = int(float64(moonchance) / float64(nbSimulations))

	result.Logs = cs.Logs

	return result
}

type Attacker struct {
	Weapon int
	Shield int
	Armour int
	ShipsInfos
}

type Defender struct {
	Weapon int
	Shield int
	Armour int
	ShipsInfos
	DefensesInfos
}

type SimulatorParams struct {
	Simulations int
}

type SimulatorResult struct {
	Simulations    int
	AttackerWin    int
	DefenderWin    int
	Draw           int
	Rounds         int
	AttackerLosses price
	DefenderLosses price
	Debris         price
	Recycler       int
	Moonchance     int
	Logs           string
}
