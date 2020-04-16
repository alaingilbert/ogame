package ogame

import (
	"errors"
	"time"
)

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
	tx               *Prioritize
	origin           Celestial
	destination      Coordinate
	speed            Speed
	mission          MissionID
	ships            ShipsInfos
	resources        Resources
	err              error
	fleet            Fleet
	minimumDeuterium int64
	expeditiontime   int64
	unionID          int64
	allShips         bool
	recallIn         int64
	successCallbacks []func(Fleet)
	errorCallbacks   []func(error)
}

// NewFleetBuilder ...
func NewFleetBuilder(b Wrapper) *FleetBuilder {
	fb := new(FleetBuilder)
	fb.b = b
	fb.mission = Transport
	fb.speed = HundredPercent
	return fb
}

// SetTx ...
func (f *FleetBuilder) SetTx(tx *Prioritize) *FleetBuilder {
	f.tx = tx
	return f
}

// SetOrigin ...
func (f *FleetBuilder) SetOrigin(v interface{}) *FleetBuilder {
	f.origin = f.b.GetCachedCelestial(v)
	return f
}

// SetDestination ...
func (f *FleetBuilder) SetDestination(v interface{}) *FleetBuilder {
	var c Celestial
	if celestial, ok := v.(Celestial); ok {
		f.destination = celestial.GetCoordinate()
	} else if planet, ok := v.(Planet); ok {
		f.destination = planet.GetCoordinate()
	} else if moon, ok := v.(Moon); ok {
		f.destination = moon.GetCoordinate()
	} else if coord, ok := v.(Coordinate); ok {
		f.destination = coord
	} else if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return f
		}
		f.destination = coord
	} else {
		c = f.b.GetCachedCelestial(v)
		if c != nil {
			f.destination = c.GetCoordinate()
		}
	}
	return f
}

// SetSpeed ...
func (f *FleetBuilder) SetSpeed(speed Speed) *FleetBuilder {
	f.speed = speed
	return f
}

// SetResources ...
func (f *FleetBuilder) SetResources(resources Resources) *FleetBuilder {
	f.resources = resources
	return f
}

// SetAllResources will send all resources from the origin
func (f *FleetBuilder) SetAllResources() *FleetBuilder {
	f.resources = Resources{Metal: -1, Crystal: -1, Deuterium: -1}
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
func (f *FleetBuilder) SetMission(mission MissionID) *FleetBuilder {
	f.mission = mission
	return f
}

// SetDuration set expedition duration
func (f *FleetBuilder) SetDuration(expeditiontime int64) *FleetBuilder {
	f.expeditiontime = expeditiontime
	return f
}

// SetUnionID set union id to join
func (f *FleetBuilder) SetUnionID(unionID int64) *FleetBuilder {
	f.unionID = unionID
	return f
}

// AddShips ...
func (f *FleetBuilder) AddShips(id ID, nbr int64) *FleetBuilder {
	f.ships.Set(id, f.ships.ByID(id)+nbr)
	return f
}

// SetShips ...
func (f *FleetBuilder) SetShips(ships ShipsInfos) *FleetBuilder {
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
	ships := f.ships
	if f.allShips {
		if f.tx != nil {
			ships, _ = f.tx.GetShips(f.origin.GetID())
		} else {
			ships, _ = f.b.GetShips(f.origin.GetID())
		}
	}
	return f.b.FlightTime(f.origin.GetCoordinate(), f.destination, f.speed, ships)
}

func (f *FleetBuilder) sendNow(tx *Prioritize) error {
	if f.origin == nil {
		f.err = errors.New("invalid origin")
		return f.err
	}

	// Set all ships
	if f.allShips {
		f.ships, _ = tx.GetShips(f.origin.GetID())
	}

	var fuel int64
	var planetResources Resources
	if f.minimumDeuterium > 0 {
		planetResources, _ = tx.GetResources(f.origin.GetID())
		_, fuel = tx.FlightTime(f.origin.GetCoordinate(), f.destination, f.speed, f.ships)
	}

	if f.minimumDeuterium > 0 && f.resources.Deuterium > 0 {
		planetResources.Deuterium = planetResources.Deuterium - (fuel + 10) - f.minimumDeuterium
		if f.resources.Deuterium > planetResources.Deuterium {
			f.resources.Deuterium = planetResources.Deuterium
		}
	}

	payload := f.resources
	// Send all resources
	if f.resources.Metal == -1 || f.resources.Crystal == -1 || f.resources.Deuterium == -1 {
		// Calculate cargo
		techs := tx.GetResearch()
		cargoCapacity := f.ships.Cargo(techs, f.b.GetServer().Settings.EspionageProbeRaids == 1, f.b.CharacterClass() == Collector)
		if f.minimumDeuterium <= 0 {
			planetResources, _ = tx.GetResources(f.origin.GetID())
		}
		if f.resources.Deuterium == -1 {
			if f.minimumDeuterium > 0 {
				planetResources.Deuterium = planetResources.Deuterium - (fuel + 10) - f.minimumDeuterium
			}
			payload.Deuterium = MinInt(cargoCapacity, planetResources.Deuterium)
			cargoCapacity -= payload.Deuterium
		}
		if f.resources.Crystal == -1 {
			payload.Crystal = MinInt(cargoCapacity, planetResources.Crystal)
			cargoCapacity -= payload.Crystal
		}
		if f.resources.Metal == -1 {
			payload.Metal = MinInt(cargoCapacity, planetResources.Metal)
		}
	}

	f.fleet, f.err = tx.EnsureFleet(f.origin.GetID(), f.ships.ToQuantifiables(), f.speed, f.destination, f.mission, payload, f.expeditiontime, f.unionID)
	return f.err
}

// SendNow send the fleet with defined configurations
func (f *FleetBuilder) SendNow() (Fleet, error) {
	var err error
	if f.tx != nil {
		err = f.sendNow(f.tx)
	} else {
		err = f.b.Tx(func(tx *Prioritize) error {
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
func (f *FleetBuilder) OnSuccess(clb func(Fleet)) *FleetBuilder {
	f.successCallbacks = append(f.successCallbacks, clb)
	return f
}
