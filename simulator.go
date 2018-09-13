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

func isAlive(unit *combatUnit) bool {
	return getUnitHull(unit) > 0
}

type combatUnit struct {
	PackedInfos uint64
}

const (
	// MAX_ARMOUR_LEVEL uint64 = 36
	// MAX_SHIELD_LEVEL uint64 = 42
	idMask     uint64 = 31
	shieldMask uint64 = 8388576
	hullMask   uint64 = 35184363700224
)

func getUnitID(unit *combatUnit) ID {
	return ID((unit.PackedInfos & idMask) >> 0)
}

func getUnitShield(unit *combatUnit) uint64 {
	return (unit.PackedInfos & shieldMask) >> 5
}

func getUnitHull(unit *combatUnit) uint64 {
	return (unit.PackedInfos & hullMask) >> 23
}

func setUnitID(unit *combatUnit, id ID) {
	unit.PackedInfos &= ^idMask
	unit.PackedInfos |= uint64(id) << 0
}

func setUnitShield(unit *combatUnit, shield uint64) {
	unit.PackedInfos &= ^shieldMask
	unit.PackedInfos |= shield << 5
}

func setUnitHull(unit *combatUnit, hull uint64) {
	unit.PackedInfos &= ^hullMask
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

func getUnitPrice(unitID ID) price {
	switch unitID {
	case SmallCargoID:
		return price{2000, 2000, 0}
	case LargeCargoID:
		return price{6000, 6000, 0}
	case LightFighterID:
		return price{3000, 1000, 0}
	case HeavyFighterID:
		return price{6000, 4000, 0}
	case CruiserID:
		return price{20000, 7000, 2000}
	case BattleshipID:
		return price{45000, 15000, 0}
	case ColonyShipID:
		return price{10000, 20000, 10000}
	case RecyclerID:
		return price{10000, 6000, 2000}
	case EspionageProbeID:
		return price{0, 1000, 0}
	case BomberID:
		return price{50000, 25000, 15000}
	case SolarSatelliteID:
		return price{0, 2000, 500}
	case DestroyerID:
		return price{60000, 50000, 15000}
	case DeathstarID:
		return price{5000000, 4000000, 1000000}
	case BattlecruiserID:
		return price{30000, 40000, 15000}
	case RocketLauncherID:
		return price{2000, 0, 0}
	case LightLaserID:
		return price{1500, 500, 0}
	case HeavyLaserID:
		return price{6000, 2000, 0}
	case GaussCannonID:
		return price{20000, 15000, 2000}
	case IonCannonID:
		return price{2000, 6000, 0}
	case PlasmaTurretID:
		return price{50000, 50000, 30000}
	case SmallShieldDomeID:
		return price{10000, 10000, 0}
	case LargeShieldDomeID:
		return price{50000, 50000, 0}
	}
	return price{0, 0, 0}
}

func getUnitBaseShield(unitID ID) int {
	switch unitID {
	case SmallCargoID:
		return 10
	case LargeCargoID:
		return 25
	case LightFighterID:
		return 10
	case HeavyFighterID:
		return 25
	case CruiserID:
		return 50
	case BattleshipID:
		return 200
	case ColonyShipID:
		return 100
	case RecyclerID:
		return 10
	case EspionageProbeID:
		return 1 // 0.01
	case BomberID:
		return 500
	case SolarSatelliteID:
		return 1
	case DestroyerID:
		return 500
	case DeathstarID:
		return 50000
	case BattlecruiserID:
		return 400
	case RocketLauncherID:
		return 20
	case LightLaserID:
		return 25
	case HeavyLaserID:
		return 100
	case GaussCannonID:
		return 200
	case IonCannonID:
		return 500
	case PlasmaTurretID:
		return 300
	case SmallShieldDomeID:
		return 2000
	case LargeShieldDomeID:
		return 10000
	}
	return 0
}

func getUnitBaseWeapon(unitID ID) uint64 {
	switch unitID {
	case SmallCargoID:
		return 5
	case LargeCargoID:
		return 5
	case LightFighterID:
		return 50
	case HeavyFighterID:
		return 150
	case CruiserID:
		return 400
	case BattleshipID:
		return 1000
	case ColonyShipID:
		return 50
	case RecyclerID:
		return 1
	case EspionageProbeID:
		return 1 // 0.01
	case BomberID:
		return 1000
	case SolarSatelliteID:
		return 1
	case DestroyerID:
		return 2000
	case DeathstarID:
		return 200000
	case BattlecruiserID:
		return 700
	case RocketLauncherID:
		return 80
	case LightLaserID:
		return 100
	case HeavyLaserID:
		return 250
	case GaussCannonID:
		return 1100
	case IonCannonID:
		return 150
	case PlasmaTurretID:
		return 3000
	case SmallShieldDomeID:
		return 1
	case LargeShieldDomeID:
		return 1
	}
	return 0
}

func getRapidFireAgainst(unit *combatUnit, targetUnit *combatUnit) int {
	rf := 0
	switch getUnitID(unit) {
	case SmallCargoID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case LargeCargoID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case LightFighterID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case HeavyFighterID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		case SmallCargoID:
			rf = 3
			break
		}
		break
	case CruiserID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		case LightFighterID:
			rf = 6
			break
		case RocketLauncherID:
			rf = 10
			break
		}
		break
	case BattleshipID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case ColonyShipID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case RecyclerID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		}
		break
	case BomberID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		case IonCannonID:
			rf = 10
			break
		case RocketLauncherID:
			rf = 20
			break
		case LightLaserID:
			rf = 20
			break
		case HeavyLaserID:
			rf = 10
			break
		}
		break
	case DestroyerID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		case LightLaserID:
			rf = 10
			break
		case BattlecruiserID:
			rf = 2
			break
		}
		break
	case DeathstarID:
		switch getUnitID(targetUnit) {
		case SmallCargoID:
			rf = 250
			break
		case LargeCargoID:
			rf = 250
			break
		case LightFighterID:
			rf = 200
			break
		case HeavyFighterID:
			rf = 100
			break
		case CruiserID:
			rf = 33
			break
		case BattleshipID:
			rf = 30
			break
		case ColonyShipID:
			rf = 250
			break
		case RecyclerID:
			rf = 250
			break
		case EspionageProbeID:
			rf = 1250
			break
		case SolarSatelliteID:
			rf = 1250
			break
		case BomberID:
			rf = 25
			break
		case DestroyerID:
			rf = 5
			break
		case RocketLauncherID:
			rf = 200
			break
		case LightLaserID:
			rf = 200
			break
		case HeavyLaserID:
			rf = 100
			break
		case GaussCannonID:
			rf = 50
			break
		case IonCannonID:
			rf = 100
			break
		case BattlecruiserID:
			rf = 15
			break
		}
		break
	case BattlecruiserID:
		switch getUnitID(targetUnit) {
		case EspionageProbeID:
			rf = 5
			break
		case SolarSatelliteID:
			rf = 5
			break
		case SmallCargoID:
			rf = 3
			break
		case LargeCargoID:
			rf = 3
			break
		case HeavyFighterID:
			rf = 4
			break
		case CruiserID:
			rf = 4
			break
		case BattleshipID:
			rf = 7
			break
		}
		break
	}
	return rf
}

func getUnitName(unitID ID) string {
	switch unitID {
	case SmallCargoID:
		return "Small cargo"
	case LargeCargoID:
		return "Large cargo"
	case LightFighterID:
		return "Light fighter"
	case HeavyFighterID:
		return "Heavy fighter"
	case CruiserID:
		return "Cruiser"
	case BattleshipID:
		return "Battleship"
	case ColonyShipID:
		return "Colony ship"
	case RecyclerID:
		return "Recycler"
	case EspionageProbeID:
		return "Expionage probe"
	case BomberID:
		return "Bomber"
	case SolarSatelliteID:
		return "Solar satellite"
	case DestroyerID:
		return "Destroyer"
	case DeathstarID:
		return "Deathstar"
	case BattlecruiserID:
		return "Battlecruiser"
	case RocketLauncherID:
		return "Rocket launcher"
	case LightLaserID:
		return "Light laser"
	case HeavyLaserID:
		return "Heavy laser"
	case GaussCannonID:
		return "Gauss cannon"
	case IonCannonID:
		return "Ion cannon"
	case PlasmaTurretID:
		return "Plasma turret"
	case SmallShieldDomeID:
		return "Small shield dome"
	case LargeShieldDomeID:
		return "Large shield dome"
	}
	return ""
}

func getUnitWeaponPower(unitID ID, weaponTechno int) uint64 {
	return uint64(float64(getUnitBaseWeapon(unitID)) * (1 + 0.1*float64(weaponTechno)))
}

func getUnitInitialShield(unitID ID, shieldTechno int) uint64 {
	return uint64(float64(getUnitBaseShield(unitID)) * (1 + 0.1*float64(shieldTechno)))
}

func getUnitInitialHullPlating(armourTechno, metalPrice, crystalPrice int) uint64 {
	return uint64((1 + (float64(armourTechno) / 10)) * (float64(metalPrice+crystalPrice) / 10))
}

func newUnit(entity *entity, unitID ID) combatUnit {
	var unit combatUnit
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
	Units           []combatUnit
	Losses          price
}

func (e *entity) Init() {
	e.Reset()
	idx := 0
	for i := 0; i < e.SmallCargo; i++ {
		e.Units[idx] = newUnit(e, SmallCargoID)
		idx++
	}
	for i := 0; i < e.LargeCargo; i++ {
		e.Units[idx] = newUnit(e, LargeCargoID)
		idx++
	}
	for i := 0; i < e.LightFighter; i++ {
		e.Units[idx] = newUnit(e, LightFighterID)
		idx++
	}
	for i := 0; i < e.HeavyFighter; i++ {
		e.Units[idx] = newUnit(e, HeavyFighterID)
		idx++
	}
	for i := 0; i < e.Cruiser; i++ {
		e.Units[idx] = newUnit(e, CruiserID)
		idx++
	}
	for i := 0; i < e.Battleship; i++ {
		e.Units[idx] = newUnit(e, BattleshipID)
		idx++
	}
	for i := 0; i < e.ColonyShip; i++ {
		e.Units[idx] = newUnit(e, ColonyShipID)
		idx++
	}
	for i := 0; i < e.Recycler; i++ {
		e.Units[idx] = newUnit(e, RecyclerID)
		idx++
	}
	for i := 0; i < e.EspionageProbe; i++ {
		e.Units[idx] = newUnit(e, EspionageProbeID)
		idx++
	}
	for i := 0; i < e.Bomber; i++ {
		e.Units[idx] = newUnit(e, BomberID)
		idx++
	}
	for i := 0; i < e.SolarSatellite; i++ {
		e.Units[idx] = newUnit(e, SolarSatelliteID)
		idx++
	}
	for i := 0; i < e.Destroyer; i++ {
		e.Units[idx] = newUnit(e, DestroyerID)
		idx++
	}
	for i := 0; i < e.Deathstar; i++ {
		e.Units[idx] = newUnit(e, DeathstarID)
		idx++
	}
	for i := 0; i < e.Battlecruiser; i++ {
		e.Units[idx] = newUnit(e, BattlecruiserID)
		idx++
	}
	for i := 0; i < e.RocketLauncher; i++ {
		e.Units[idx] = newUnit(e, RocketLauncherID)
		idx++
	}
	for i := 0; i < e.LightLaser; i++ {
		e.Units[idx] = newUnit(e, LightLaserID)
		idx++
	}
	for i := 0; i < e.HeavyLaser; i++ {
		e.Units[idx] = newUnit(e, HeavyLaserID)
		idx++
	}
	for i := 0; i < e.GaussCannon; i++ {
		e.Units[idx] = newUnit(e, GaussCannonID)
		idx++
	}
	for i := 0; i < e.IonCannon; i++ {
		e.Units[idx] = newUnit(e, IonCannonID)
		idx++
	}
	for i := 0; i < e.PlasmaTurret; i++ {
		e.Units[idx] = newUnit(e, PlasmaTurretID)
		idx++
	}
	for i := 0; i < e.SmallShieldDome; i++ {
		e.Units[idx] = newUnit(e, SmallShieldDomeID)
		idx++
	}
	for i := 0; i < e.LargeShieldDome; i++ {
		e.Units[idx] = newUnit(e, LargeShieldDomeID)
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

func (simulator *combatSimulator) hasExploded(entity *entity, defendingUnit *combatUnit) bool {
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

func (simulator *combatSimulator) getAnotherShot(unit, targetUnit *combatUnit) bool {
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

func (simulator *combatSimulator) attack(attacker *entity, attackingUnit *combatUnit, defender *entity, defendingUnit *combatUnit) {
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

func isShip(unit *combatUnit) bool {
	return getUnitID(unit).IsShip()
}

func (simulator *combatSimulator) RemoveDestroyedUnits() {
	l := simulator.Defender.TotalUnits
	for i := l - 1; i >= 0; i-- {
		unit := &simulator.Defender.Units[i]
		if getUnitHull(unit) <= 0 {
			unitPrice := getUnitPrice(getUnitID(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Defender.Losses.Add(unitPrice)
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
		if getUnitHull(unit) <= 0 {
			unitPrice := getUnitPrice(getUnitID(unit))
			if isShip(unit) {
				simulator.Debris.Metal += int(simulator.FleetToDebris * float64(unitPrice.Metal))
				simulator.Debris.Crystal += int(simulator.FleetToDebris * float64(unitPrice.Crystal))
			}
			simulator.Attacker.Losses.Add(unitPrice)
			simulator.Attacker.Units[i] = simulator.Attacker.Units[simulator.Attacker.TotalUnits-1]
			simulator.Attacker.TotalUnits--
			//simulator.Attacker.Units = simulator.Attacker.Units[:len(simulator.Attacker.Units)-1]
			if simulator.IsLogging {
				simulator.Logs += fmt.Sprintf("%s lost all its integrity, remove from battle\n", getUnitName(getUnitID(unit)))
			}
		}
	}
}

func (simulator *combatSimulator) RestoreShields() {
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
	attacker.Units = make([]combatUnit, attacker.TotalUnits+1, attacker.TotalUnits+1)

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
	defender.Units = make([]combatUnit, defender.TotalUnits+1, defender.TotalUnits+1)

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

// Attacker ...
type Attacker struct {
	Weapon int
	Shield int
	Armour int
	ShipsInfos
}

// Defender ...
type Defender struct {
	Weapon int
	Shield int
	Armour int
	ShipsInfos
	DefensesInfos
}

// SimulatorParams ...
type SimulatorParams struct {
	Simulations int
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
