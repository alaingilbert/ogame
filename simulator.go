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

const (
	smallCargoConst = iota
	largeCargoConst
	lightFighterConst
	heavyFighterConst
	cruiserConst
	battleshipConst
	colonyShipConst
	recyclerConst
	espionageProbeConst
	bomberConst
	solarSatelliteConst
	destroyerConst
	deathstarConst
	battlecruiserConst
	rocketLauncherConst
	lightLaserConst
	heavyLaserConst
	gaussCannonConst
	ionCannonConst
	plasmaTurretConst
	smallShieldDomeConst
	largeShieldDomeConst
)

func isAlive(unit *CombatUnit) bool {
	return getUnitHull(unit) > 0
}

// CombatUnit ...
type CombatUnit struct {
	PackedInfos uint64
}

const (
	maxArmourLevel uint64 = 36
	maxShieldLevel uint64 = 42
	idMask         uint64 = 31
	shieldMask     uint64 = 8388576
	hullMask       uint64 = 35184363700224
)

func getUnitID(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & idMask) >> 0
}

func getUnitShield(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & shieldMask) >> 5
}

func getUnitHull(unit *CombatUnit) uint64 {
	return (unit.PackedInfos & hullMask) >> 23
}

func setUnitID(unit *CombatUnit, id uint64) {
	unit.PackedInfos &= ^idMask
	unit.PackedInfos |= id << 0
}

func setUnitShield(unit *CombatUnit, shield uint64) {
	unit.PackedInfos &= ^shieldMask
	unit.PackedInfos |= shield << 5
}

func setUnitHull(unit *CombatUnit, hull uint64) {
	unit.PackedInfos &= ^hullMask
	unit.PackedInfos |= hull << 23
}

type price struct {
	Metal     int
	Crystal   int
	Deuterium int
}

func (p price) Total() int {
	return p.Metal + p.Crystal + p.Deuterium
}

func (p *price) add(n price) {
	p.Metal += n.Metal
	p.Crystal += n.Crystal
	p.Deuterium += n.Deuterium
}

func getUnitPrice(unitID uint64) price {
	switch unitID {
	case smallCargoConst:
		return price{2000, 2000, 0}
	case largeCargoConst:
		return price{6000, 6000, 0}
	case lightFighterConst:
		return price{3000, 1000, 0}
	case heavyFighterConst:
		return price{6000, 4000, 0}
	case cruiserConst:
		return price{20000, 7000, 2000}
	case battleshipConst:
		return price{45000, 15000, 0}
	case colonyShipConst:
		return price{10000, 20000, 10000}
	case recyclerConst:
		return price{10000, 6000, 2000}
	case espionageProbeConst:
		return price{0, 1000, 0}
	case bomberConst:
		return price{50000, 25000, 15000}
	case solarSatelliteConst:
		return price{0, 2000, 500}
	case destroyerConst:
		return price{60000, 50000, 15000}
	case deathstarConst:
		return price{5000000, 4000000, 1000000}
	case battlecruiserConst:
		return price{30000, 40000, 15000}
	case rocketLauncherConst:
		return price{2000, 0, 0}
	case lightLaserConst:
		return price{1500, 500, 0}
	case heavyLaserConst:
		return price{6000, 2000, 0}
	case gaussCannonConst:
		return price{20000, 15000, 2000}
	case ionCannonConst:
		return price{2000, 6000, 0}
	case plasmaTurretConst:
		return price{50000, 50000, 30000}
	case smallShieldDomeConst:
		return price{10000, 10000, 0}
	case largeShieldDomeConst:
		return price{50000, 50000, 0}
	}
	return price{0, 0, 0}
}

func getUnitBaseShield(unitID uint64) int {
	switch unitID {
	case smallCargoConst:
		return 10
	case largeCargoConst:
		return 25
	case lightFighterConst:
		return 10
	case heavyFighterConst:
		return 25
	case cruiserConst:
		return 50
	case battleshipConst:
		return 200
	case colonyShipConst:
		return 100
	case recyclerConst:
		return 10
	case espionageProbeConst:
		return 1 // 0.01
	case bomberConst:
		return 500
	case solarSatelliteConst:
		return 1
	case destroyerConst:
		return 500
	case deathstarConst:
		return 50000
	case battlecruiserConst:
		return 400
	case rocketLauncherConst:
		return 20
	case lightLaserConst:
		return 25
	case heavyLaserConst:
		return 100
	case gaussCannonConst:
		return 200
	case ionCannonConst:
		return 500
	case plasmaTurretConst:
		return 300
	case smallShieldDomeConst:
		return 2000
	case largeShieldDomeConst:
		return 10000
	}
	return 0
}

func getUnitBaseWeapon(unitID uint64) uint64 {
	switch unitID {
	case smallCargoConst:
		return 5
	case largeCargoConst:
		return 5
	case lightFighterConst:
		return 50
	case heavyFighterConst:
		return 150
	case cruiserConst:
		return 400
	case battleshipConst:
		return 1000
	case colonyShipConst:
		return 50
	case recyclerConst:
		return 1
	case espionageProbeConst:
		return 1 // 0.01
	case bomberConst:
		return 1000
	case solarSatelliteConst:
		return 1
	case destroyerConst:
		return 2000
	case deathstarConst:
		return 200000
	case battlecruiserConst:
		return 700
	case rocketLauncherConst:
		return 80
	case lightLaserConst:
		return 100
	case heavyLaserConst:
		return 250
	case gaussCannonConst:
		return 1100
	case ionCannonConst:
		return 150
	case plasmaTurretConst:
		return 3000
	case smallShieldDomeConst:
		return 1
	case largeShieldDomeConst:
		return 1
	}
	return 0
}

func getRapidFireAgainst(unit *CombatUnit, targetUnit *CombatUnit) int {
	rf := 0
	switch getUnitID(unit) {
	case smallCargoConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case largeCargoConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case lightFighterConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case heavyFighterConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		case smallCargoConst:
			rf = 3
		}
	case cruiserConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		case lightFighterConst:
			rf = 6
		case rocketLauncherConst:
			rf = 10
		}
	case battleshipConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case colonyShipConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case recyclerConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		}
	case bomberConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		case ionCannonConst:
			rf = 10
		case rocketLauncherConst:
			rf = 20
		case lightLaserConst:
			rf = 20
		case heavyLaserConst:
			rf = 10
		}
	case destroyerConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		case lightLaserConst:
			rf = 10
		case battlecruiserConst:
			rf = 2
		}
	case deathstarConst:
		switch getUnitID(targetUnit) {
		case smallCargoConst:
			rf = 250
		case largeCargoConst:
			rf = 250
		case lightFighterConst:
			rf = 200
		case heavyFighterConst:
			rf = 100
		case cruiserConst:
			rf = 33
		case battleshipConst:
			rf = 30
		case colonyShipConst:
			rf = 250
		case recyclerConst:
			rf = 250
		case espionageProbeConst:
			rf = 1250
		case solarSatelliteConst:
			rf = 1250
		case bomberConst:
			rf = 25
		case destroyerConst:
			rf = 5
		case rocketLauncherConst:
			rf = 200
		case lightLaserConst:
			rf = 200
		case heavyLaserConst:
			rf = 100
		case gaussCannonConst:
			rf = 50
		case ionCannonConst:
			rf = 100
		case battlecruiserConst:
			rf = 15
		}
	case battlecruiserConst:
		switch getUnitID(targetUnit) {
		case espionageProbeConst:
			rf = 5
		case solarSatelliteConst:
			rf = 5
		case smallCargoConst:
			rf = 3
		case largeCargoConst:
			rf = 3
		case heavyFighterConst:
			rf = 4
		case cruiserConst:
			rf = 4
		case battleshipConst:
			rf = 7
		}
	}
	return rf
}

func getUnitName(unitID uint64) string {
	switch unitID {
	case smallCargoConst:
		return "Small cargo"
	case largeCargoConst:
		return "Large cargo"
	case lightFighterConst:
		return "Light fighter"
	case heavyFighterConst:
		return "Heavy fighter"
	case cruiserConst:
		return "Cruiser"
	case battleshipConst:
		return "Battleship"
	case colonyShipConst:
		return "Colony ship"
	case recyclerConst:
		return "Recycler"
	case espionageProbeConst:
		return "Expionage probe"
	case bomberConst:
		return "Bomber"
	case solarSatelliteConst:
		return "Solar satellite"
	case destroyerConst:
		return "Destroyer"
	case deathstarConst:
		return "Deathstar"
	case battlecruiserConst:
		return "Battlecruiser"
	case rocketLauncherConst:
		return "Rocket launcher"
	case lightLaserConst:
		return "Light laser"
	case heavyLaserConst:
		return "Heavy laser"
	case gaussCannonConst:
		return "Gauss cannon"
	case ionCannonConst:
		return "Ion cannon"
	case plasmaTurretConst:
		return "Plasma turret"
	case smallShieldDomeConst:
		return "Small shield dome"
	case largeShieldDomeConst:
		return "Large shield dome"
	}
	return ""
}

func getUnitWeaponPower(unitID uint64, weaponTechno int) uint64 {
	return uint64(float64(getUnitBaseWeapon(unitID)) * (1 + 0.1*float64(weaponTechno)))
}

func getUnitInitialShield(unitID uint64, shieldTechno int) uint64 {
	return uint64(float64(getUnitBaseShield(unitID)) * (1 + 0.1*float64(shieldTechno)))
}

func getUnitInitialHullPlating(armourTechno, metalPrice, crystalPrice int) uint64 {
	return uint64((1 + (float64(armourTechno) / 10)) * (float64(metalPrice+crystalPrice) / 10))
}

func newUnit(entity *entity, unitID uint64) CombatUnit {
	var unit CombatUnit
	setUnitID(&unit, unitID)
	unitPrice := getUnitPrice(unitID)
	setUnitHull(&unit, getUnitInitialHullPlating(entity.Armour, unitPrice.Metal, unitPrice.Crystal))
	setUnitShield(&unit, getUnitInitialShield(unitID, entity.Shield))
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

func (e *entity) init() {
	e.reset()
	idx := 0
	for i := 0; i < e.SmallCargo; i++ {
		e.Units[idx] = newUnit(e, smallCargoConst)
		idx++
	}
	for i := 0; i < e.LargeCargo; i++ {
		e.Units[idx] = newUnit(e, largeCargoConst)
		idx++
	}
	for i := 0; i < e.LightFighter; i++ {
		e.Units[idx] = newUnit(e, lightFighterConst)
		idx++
	}
	for i := 0; i < e.HeavyFighter; i++ {
		e.Units[idx] = newUnit(e, heavyFighterConst)
		idx++
	}
	for i := 0; i < e.Cruiser; i++ {
		e.Units[idx] = newUnit(e, cruiserConst)
		idx++
	}
	for i := 0; i < e.Battleship; i++ {
		e.Units[idx] = newUnit(e, battleshipConst)
		idx++
	}
	for i := 0; i < e.ColonyShip; i++ {
		e.Units[idx] = newUnit(e, colonyShipConst)
		idx++
	}
	for i := 0; i < e.Recycler; i++ {
		e.Units[idx] = newUnit(e, recyclerConst)
		idx++
	}
	for i := 0; i < e.EspionageProbe; i++ {
		e.Units[idx] = newUnit(e, espionageProbeConst)
		idx++
	}
	for i := 0; i < e.Bomber; i++ {
		e.Units[idx] = newUnit(e, bomberConst)
		idx++
	}
	for i := 0; i < e.SolarSatellite; i++ {
		e.Units[idx] = newUnit(e, solarSatelliteConst)
		idx++
	}
	for i := 0; i < e.Destroyer; i++ {
		e.Units[idx] = newUnit(e, destroyerConst)
		idx++
	}
	for i := 0; i < e.Deathstar; i++ {
		e.Units[idx] = newUnit(e, deathstarConst)
		idx++
	}
	for i := 0; i < e.Battlecruiser; i++ {
		e.Units[idx] = newUnit(e, battlecruiserConst)
		idx++
	}
	for i := 0; i < e.RocketLauncher; i++ {
		e.Units[idx] = newUnit(e, rocketLauncherConst)
		idx++
	}
	for i := 0; i < e.LightLaser; i++ {
		e.Units[idx] = newUnit(e, lightLaserConst)
		idx++
	}
	for i := 0; i < e.HeavyLaser; i++ {
		e.Units[idx] = newUnit(e, heavyLaserConst)
		idx++
	}
	for i := 0; i < e.GaussCannon; i++ {
		e.Units[idx] = newUnit(e, gaussCannonConst)
		idx++
	}
	for i := 0; i < e.IonCannon; i++ {
		e.Units[idx] = newUnit(e, ionCannonConst)
		idx++
	}
	for i := 0; i < e.PlasmaTurret; i++ {
		e.Units[idx] = newUnit(e, plasmaTurretConst)
		idx++
	}
	for i := 0; i < e.SmallShieldDome; i++ {
		e.Units[idx] = newUnit(e, smallShieldDomeConst)
		idx++
	}
	for i := 0; i < e.LargeShieldDome; i++ {
		e.Units[idx] = newUnit(e, largeShieldDomeConst)
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
	unitPrice := getUnitPrice(getUnitID(defendingUnit))
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
				msg += fmt.Sprintf("%s gets another shot.", getUnitName(getUnitID(unit)))
			}
		} else {
			if simulator.IsLogging {
				msg += fmt.Sprintf("%s does not get another shot.", getUnitName(getUnitID(unit)))
			}
			rapidFire = false
		}
	} else {
		if simulator.IsLogging {
			msg += fmt.Sprintf("%s doesn't have rapid fire against %s.", getUnitName(getUnitID(unit)), getUnitName(getUnitID(targetUnit)))
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
		simulator.Logs += fmt.Sprintf("%s fires at %s; ", getUnitName(getUnitID(attackingUnit)), getUnitName(getUnitID(defendingUnit)))
	}

	weapon := getUnitWeaponPower(getUnitID(attackingUnit), attacker.Weapon)
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
		simulator.Logs += fmt.Sprintf("result is %s %d %d\n", getUnitName(getUnitID(defendingUnit)), getUnitHull(defendingUnit), getUnitShield(defendingUnit))
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
	switch getUnitID(unit) {
	case smallCargoConst:
		return true
	case largeCargoConst:
		return true
	case lightFighterConst:
		return true
	case heavyFighterConst:
		return true
	case cruiserConst:
		return true
	case battleshipConst:
		return true
	case colonyShipConst:
		return true
	case recyclerConst:
		return true
	case espionageProbeConst:
		return true
	case bomberConst:
		return true
	case solarSatelliteConst:
		return true
	case destroyerConst:
		return true
	case deathstarConst:
		return true
	case battlecruiserConst:
		return true
	}
	return false
}

func (simulator *combatSimulator) removeDestroyedUnits() {
	l := simulator.Defender.TotalUnits
	for i := l - 1; i >= 0; i-- {
		unit := &simulator.Defender.Units[i]
		if getUnitHull(unit) == 0 {
			unitPrice := getUnitPrice(getUnitID(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Defender.Losses.add(unitPrice)
			simulator.Defender.Units[i] = simulator.Defender.Units[simulator.Defender.TotalUnits-1]
			simulator.Defender.TotalUnits--
			//simulator.Defender.Units = simulator.Defender.Units[:len(simulator.Defender.Units)-1]
			if simulator.IsLogging {
				simulator.Logs += fmt.Sprintf("%s lost all its integrity, remove from battle\n", getUnitName(getUnitID(unit)))
			}
		}
	}
	l = simulator.Attacker.TotalUnits
	for i := l - 1; i >= 0; i-- {
		unit := &simulator.Attacker.Units[i]
		if getUnitHull(unit) == 0 {
			unitPrice := getUnitPrice(getUnitID(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Attacker.Losses.add(unitPrice)
			simulator.Attacker.Units[i] = simulator.Attacker.Units[simulator.Attacker.TotalUnits-1]
			simulator.Attacker.TotalUnits--
			//simulator.Attacker.Units = simulator.Attacker.Units[:len(simulator.Attacker.Units)-1]
			if simulator.IsLogging {
				simulator.Logs += fmt.Sprintf("%s lost all its integrity, remove from battle\n", getUnitName(getUnitID(unit)))
			}
		}
	}
}

func (simulator *combatSimulator) restoreShields() {
	for i := 0; i < simulator.Attacker.TotalUnits; i++ {
		unit := &simulator.Attacker.Units[i]
		setUnitShield(unit, getUnitInitialShield(getUnitID(unit), simulator.Attacker.Shield))
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("%s still has integrity, restore its shield\n", getUnitName(getUnitID(unit)))
		}
	}
	for i := 0; i < simulator.Defender.TotalUnits; i++ {
		unit := &simulator.Defender.Units[i]
		setUnitShield(unit, getUnitInitialShield(getUnitID(unit), simulator.Defender.Shield))
		if simulator.IsLogging {
			simulator.Logs += fmt.Sprintf("%s still has integrity, restore its shield\n", getUnitName(getUnitID(unit)))
		}
	}
}

func (simulator *combatSimulator) isCombatDone() bool {
	return simulator.Attacker.TotalUnits <= 0 || simulator.Defender.TotalUnits <= 0
}

func (simulator *combatSimulator) getMoonchance() int {
	debris := float64(simulator.Debris.Metal) + float64(simulator.Debris.Crystal)
	return int(math.Min(debris/100000.0, 20.0))
}

func (simulator *combatSimulator) printWinner() {
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
	simulator.Attacker.init()
	simulator.Defender.init()
	for currentRound := 1; currentRound <= simulator.MaxRounds; currentRound++ {
		simulator.Rounds = currentRound
		if simulator.IsLogging {
			simulator.Logs += strings.Repeat("-", 80) + "\n"
			simulator.Logs += "ROUND " + strconv.Itoa(currentRound) + "\n"
			simulator.Logs += strings.Repeat("-", 80) + "\n"
		}
		simulator.attackerFires()
		simulator.defenderFires()
		simulator.removeDestroyedUnits()
		simulator.restoreShields()
		if simulator.isCombatDone() {
			break
		}
	}
	simulator.printWinner()
}

func newCombatSimulator(attacker *entity, defender *entity) *combatSimulator {
	cs := new(combatSimulator)
	cs.Attacker = *attacker
	cs.Defender = *defender
	cs.IsLogging = false
	cs.MaxRounds = 6
	return cs
}

// Config ...
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

func (e *entity) reset() {
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

// Simulate ...
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
	attacker.SmallCargo = int(attackerParam.SmallCargo)
	attacker.LargeCargo = int(attackerParam.LargeCargo)
	attacker.LightFighter = int(attackerParam.LightFighter)
	attacker.HeavyFighter = int(attackerParam.HeavyFighter)
	attacker.Cruiser = int(attackerParam.Cruiser)
	attacker.Battleship = int(attackerParam.Battleship)
	attacker.ColonyShip = int(attackerParam.ColonyShip)
	attacker.Recycler = int(attackerParam.Recycler)
	attacker.EspionageProbe = int(attackerParam.EspionageProbe)
	attacker.Bomber = int(attackerParam.Bomber)
	attacker.SolarSatellite = 0
	attacker.Destroyer = int(attackerParam.Destroyer)
	attacker.Deathstar = int(attackerParam.Deathstar)
	attacker.Battlecruiser = int(attackerParam.Battlecruiser)
	attacker.RocketLauncher = 0
	attacker.LightLaser = 0
	attacker.HeavyLaser = 0
	attacker.GaussCannon = 0
	attacker.IonCannon = 0
	attacker.PlasmaTurret = 0
	attacker.SmallShieldDome = 0
	attacker.LargeShieldDome = 0
	attacker.reset()
	attacker.Units = make([]CombatUnit, attacker.TotalUnits+1)

	defender := newEntity()
	defender.Weapon = defenderParam.Weapon
	defender.Shield = defenderParam.Shield
	defender.Armour = defenderParam.Armour
	defender.SmallCargo = int(defenderParam.SmallCargo)
	defender.LargeCargo = int(defenderParam.LargeCargo)
	defender.LightFighter = int(defenderParam.LightFighter)
	defender.HeavyFighter = int(defenderParam.HeavyFighter)
	defender.Cruiser = int(defenderParam.Cruiser)
	defender.Battleship = int(defenderParam.Battleship)
	defender.ColonyShip = int(defenderParam.ColonyShip)
	defender.Recycler = int(defenderParam.Recycler)
	defender.EspionageProbe = int(defenderParam.EspionageProbe)
	defender.Bomber = int(defenderParam.Bomber)
	defender.SolarSatellite = int(defenderParam.SolarSatellite)
	defender.Destroyer = int(defenderParam.Destroyer)
	defender.Deathstar = int(defenderParam.Deathstar)
	defender.Battlecruiser = int(defenderParam.Battlecruiser)
	defender.RocketLauncher = int(defenderParam.RocketLauncher)
	defender.LightLaser = int(defenderParam.LightLaser)
	defender.HeavyLaser = int(defenderParam.HeavyLaser)
	defender.GaussCannon = int(defenderParam.GaussCannon)
	defender.IonCannon = int(defenderParam.IonCannon)
	defender.PlasmaTurret = int(defenderParam.PlasmaTurret)
	defender.SmallShieldDome = int(defenderParam.SmallShieldDome)
	defender.LargeShieldDome = int(defenderParam.LargeShieldDome)
	defender.reset()
	defender.Units = make([]CombatUnit, defender.TotalUnits+1)

	cs := newCombatSimulator(attacker, defender)
	cs.IsLogging = false
	cs.FleetToDebris = params.FleetToDebris

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
		attackerLosses.add(cs.Attacker.Losses)
		defenderLosses.add(cs.Defender.Losses)
		debris.add(cs.Debris)
		rounds += cs.Rounds
		moonchance += cs.getMoonchance()
	}

	result := SimulatorResult{}
	result.Simulations = nbSimulations
	result.AttackerWin = int(math.Round(float64(attackerWin) / float64(nbSimulations) * 100))
	result.DefenderWin = int(math.Round(float64(defenderWin) / float64(nbSimulations) * 100))
	result.Draw = int(math.Round(float64(draw) / float64(nbSimulations) * 100))
	result.Rounds = int(math.Round(float64(rounds) / float64(nbSimulations)))
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

// Attacker ...
type Attacker struct {
	Weapon int
	Shield int
	Armour int
	ShipsInfos
}

// Defender ...
type Defender struct {
	Metal     int
	Crystal   int
	Deuterium int
	Weapon    int
	Shield    int
	Armour    int
	ShipsInfos
	DefensesInfos
}

// SimulatorParams ...
type SimulatorParams struct {
	Simulations   int
	FleetToDebris float64
}

// SimulatorResult ...
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
