package wrapper

import (
	"errors"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"time"
)

var ErrInvalidOrigin = errors.New("invalid origin")

// FleetBuilderFactory ...
type FleetBuilderFactory struct {
	b Wrapper
}

// NewFleetBuilderFactory ...
func NewFleetBuilderFactory(b Wrapper) *FleetBuilderFactory {
	return &FleetBuilderFactory{b: b}
}

// NewFleet ...
func (f FleetBuilderFactory) NewFleet() *FleetBuilder {
	return NewFleetBuilder(f.b)
}

// FleetBuilder ...
type FleetBuilder struct {
	b                Wrapper
	tx               Prioritizable
	origin           Celestial
	destination      ogame.Coordinate
	speed            ogame.Speed
	mission          ogame.MissionID
	ships            ogame.ShipsInfos
	resources        ogame.Resources
	err              error
	fleet            ogame.Fleet
	minimumDeuterium int64
	holdingTime      int64
	unionID          int64
	allShips         bool
	recallIn         int64
	successCallbacks []func(ogame.Fleet)
	errorCallbacks   []func(error)
}

// NewFleetBuilder ...
func NewFleetBuilder(b Wrapper) *FleetBuilder {
	fb := new(FleetBuilder)
	fb.b = b
	fb.mission = ogame.Transport
	fb.speed = ogame.HundredPercent
	return fb
}

// SetTx ...
func (f *FleetBuilder) SetTx(tx Prioritizable) *FleetBuilder {
	f.tx = tx
	return f
}

// SetOrigin ...
func (f *FleetBuilder) SetOrigin(v IntoCelestial) *FleetBuilder {
	f.origin, _ = f.b.GetCachedCelestial(v)
	return f
}

// SetDestination ...
func (f *FleetBuilder) SetDestination(v IntoCoordinate) *FleetBuilder {
	f.destination, _ = ConvertIntoCoordinate(f.b, v)
	return f
}

// SetSpeed ...
func (f *FleetBuilder) SetSpeed(speed ogame.Speed) *FleetBuilder {
	f.speed = speed
	return f
}

// SetResources ...
func (f *FleetBuilder) SetResources(resources ogame.Resources) *FleetBuilder {
	f.resources = resources
	return f
}

// SetMetal ...
func (f *FleetBuilder) SetMetal(metal int64) *FleetBuilder {
	f.resources.Metal = max(metal, -1)
	return f
}

// SetCrystal ...
func (f *FleetBuilder) SetCrystal(crystal int64) *FleetBuilder {
	f.resources.Crystal = max(crystal, -1)
	return f
}

// SetDeuterium ...
func (f *FleetBuilder) SetDeuterium(deuterium int64) *FleetBuilder {
	f.resources.Deuterium = max(deuterium, -1)
	return f
}

// SetAllResources will send all resources from the origin
func (f *FleetBuilder) SetAllResources() *FleetBuilder {
	f.resources = ogame.Resources{Metal: -1, Crystal: -1, Deuterium: -1}
	return f
}

// SetAllMetal will send all metal from the origin
func (f *FleetBuilder) SetAllMetal() *FleetBuilder {
	f.resources.Metal = -1
	return f
}

// SetAllCrystal will send all crystal from the origin
func (f *FleetBuilder) SetAllCrystal() *FleetBuilder {
	f.resources.Crystal = -1
	return f
}

// SetAllDeuterium will send all deuterium from the origin
func (f *FleetBuilder) SetAllDeuterium() *FleetBuilder {
	f.resources.Deuterium = -1
	return f
}

// SetMinimumDeuterium set minimum deuterium to keep on celestial
func (f *FleetBuilder) SetMinimumDeuterium(minimumDeuterium int64) *FleetBuilder {
	f.minimumDeuterium = minimumDeuterium
	return f
}

// SetMission ...
func (f *FleetBuilder) SetMission(mission ogame.MissionID) *FleetBuilder {
	f.mission = mission
	return f
}

// SetDuration set expedition duration
func (f *FleetBuilder) SetDuration(holdingTime int64) *FleetBuilder {
	f.holdingTime = holdingTime
	return f
}

// SetUnionID set union id to join
func (f *FleetBuilder) SetUnionID(unionID int64) *FleetBuilder {
	f.unionID = unionID
	return f
}

// AddShips ...
func (f *FleetBuilder) AddShips(id ogame.ID, nbr int64) *FleetBuilder {
	f.ships.Set(id, f.ships.ByID(id)+nbr)
	return f
}

// SetShips ...
func (f *FleetBuilder) SetShips(ships ogame.ShipsInfos) *FleetBuilder {
	f.ships = ships
	return f
}

// SetAllShips ...
func (f *FleetBuilder) SetAllShips() *FleetBuilder {
	f.allShips = true
	return f
}

// SetRecallIn ...
func (f *FleetBuilder) SetRecallIn(secs int64) *FleetBuilder {
	f.recallIn = secs
	return f
}

// FlightTime ...
func (f *FleetBuilder) FlightTime() (secs, fuel int64) {
	return f.flightTime(true)
}

// FastFlightTime ...
func (f *FleetBuilder) FastFlightTime() (secs, fuel int64) {
	return f.flightTime(false)
}

// Cargo get the fleet cargo capacity
func (f *FleetBuilder) Cargo() int64 {
	return f.cargo()
}

// CargoUsing get the fleet cargo capacity using the provided techs & lfBonuses
func (f *FleetBuilder) CargoUsing(techs ogame.IResearches, lfBonuses ogame.LfBonuses) int64 {
	return f.cargoUsing(techs, lfBonuses)
}

func (f *FleetBuilder) cargo() int64 {
	var w Prioritizable = f.b
	if f.tx != nil {
		w = f.tx
	}
	techs := w.GetCachedResearch()
	lfBonuses, _ := w.GetCachedLfBonuses()
	return f.cargoUsing(techs, lfBonuses)
}

func (f *FleetBuilder) cargoUsing(techs ogame.IResearches, lfBonuses ogame.LfBonuses) int64 {
	multiplier := float64(f.b.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	return f.ships.Cargo(techs, lfBonuses, f.b.CharacterClass(), multiplier, f.b.GetServer().OGameSettings().ProbeRaidsEnabled())
}

func (f *FleetBuilder) flightTime(checkTarget bool) (secs, fuel int64) {
	var w Prioritizable = f.b
	if f.tx != nil {
		w = f.tx
	}
	origin := f.origin
	if origin == nil {
		f.err = ErrInvalidOrigin
		return 0, 0
	}
	ships := f.ships
	if f.allShips {
		ships, _ = w.GetShips(origin.GetID())
	}
	flightTimeFn := utils.Ternary(checkTarget, w.FlightTime, w.FastFlightTime)
	return flightTimeFn(origin.GetCoordinate(), f.destination, f.speed, ships, f.mission, f.holdingTime)
}

func (f *FleetBuilder) sendNow(tx Prioritizable) error {
	if f.origin == nil {
		f.err = ErrInvalidOrigin
		return f.err
	}

	// Set all ships
	if f.allShips {
		f.ships, _ = tx.GetShips(f.origin.GetID())
	}

	var fuel int64
	var planetResources ogame.Resources
	if f.minimumDeuterium > 0 {
		planetResources, _ = tx.GetResources(f.origin.GetID())
		_, fuel = tx.FlightTime(f.origin.GetCoordinate(), f.destination, f.speed, f.ships, f.mission, f.holdingTime)
	}

	if f.minimumDeuterium > 0 && f.resources.Deuterium > 0 {
		planetResources.Deuterium = planetResources.Deuterium - (fuel + 10) - f.minimumDeuterium
		f.resources.Deuterium = min(f.resources.Deuterium, planetResources.Deuterium)
	}

	payload := f.resources
	// Send all resources
	if f.resources.Metal == -1 || f.resources.Crystal == -1 || f.resources.Deuterium == -1 {
		// Calculate cargo
		cargoCapacity := f.cargo()
		if f.minimumDeuterium <= 0 {
			planetResources, _ = tx.GetResources(f.origin.GetID())
		}
		if f.resources.Deuterium == -1 {
			if f.minimumDeuterium > 0 {
				planetResources.Deuterium = planetResources.Deuterium - (fuel + 10) - f.minimumDeuterium
			}
			payload.Deuterium = min(cargoCapacity, planetResources.Deuterium)
			cargoCapacity -= payload.Deuterium
		}
		if f.resources.Crystal == -1 {
			payload.Crystal = min(cargoCapacity, planetResources.Crystal)
			cargoCapacity -= payload.Crystal
		}
		if f.resources.Metal == -1 {
			payload.Metal = min(cargoCapacity, planetResources.Metal)
		}
	}

	f.fleet, f.err = tx.EnsureFleet(f.origin.GetID(), f.ships, f.speed, f.destination, f.mission, payload, f.holdingTime, f.unionID)
	return f.err
}

// SendNow send the fleet with defined configurations
func (f *FleetBuilder) SendNow() (ogame.Fleet, error) {
	var err error
	if f.tx != nil {
		err = f.sendNow(f.tx)
	} else {
		err = f.b.TxNamed("FleetBuilder SendNow", func(tx Prioritizable) error {
			return f.sendNow(tx)
		})
	}
	if err != nil {
		// On error, call error callbacks
		for _, clb := range f.errorCallbacks {
			clb(err)
		}
	} else {
		if f.recallIn > 0 {
			go func() {
				time.Sleep(time.Duration(f.recallIn) * time.Second)
				_ = f.b.CancelFleet(f.fleet.ID)
			}()
		}

		// Otherwise, call success callbacks
		for _, clb := range f.successCallbacks {
			clb(f.fleet)
		}
	}
	return f.fleet, f.err
}

// OnError register an error callback
func (f *FleetBuilder) OnError(clb func(error)) *FleetBuilder {
	f.errorCallbacks = append(f.errorCallbacks, clb)
	return f
}

// OnSuccess register a success callback
func (f *FleetBuilder) OnSuccess(clb func(ogame.Fleet)) *FleetBuilder {
	f.successCallbacks = append(f.successCallbacks, clb)
	return f
}
