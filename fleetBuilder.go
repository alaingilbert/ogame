package ogame

import (
	"math"
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
	origin           Celestial
	destination      Coordinate
	speed            Speed
	mission          MissionID
	ships            ShipsInfos
	resources        Resources
	err              error
	fleet            Fleet
	expeditiontime   int
	unionID          int
	allShips         bool
	recallIn         int
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

// SetMission ...
func (f *FleetBuilder) SetMission(mission MissionID) *FleetBuilder {
	f.mission = mission
	return f
}

// SetDuration set expedition duration
func (f *FleetBuilder) SetDuration(expeditiontime int) *FleetBuilder {
	f.expeditiontime = expeditiontime
	return f
}

// SetUnionID set union id to join
func (f *FleetBuilder) SetUnionID(unionID int) *FleetBuilder {
	f.unionID = unionID
	return f
}

// AddShips ...
func (f *FleetBuilder) AddShips(id ID, nbr int) *FleetBuilder {
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
func (f *FleetBuilder) SetRecallIn(secs int) *FleetBuilder {
	f.recallIn = secs
	return f
}

// FlightTime ...
func (f *FleetBuilder) FlightTime() (secs, fuel int) {
	ships := f.ships
	if f.allShips {
		ships, _ = f.b.GetShips(f.origin.GetID())
	}
	return f.b.FlightTime(f.origin.GetCoordinate(), f.destination, f.speed, ships)
}

// SendNow send the fleet with defined configurations
func (f *FleetBuilder) SendNow() (Fleet, error) {
	err := f.b.Tx(func(tx *Prioritize) error {

		// Set all ships
		if f.allShips {
			f.ships, _ = tx.GetShips(f.origin.GetID())
		}

		payload := f.resources
		// Send all resources
		if f.resources.Metal == -1 && f.resources.Crystal == -1 && f.resources.Deuterium == -1 {
			// Calculate cargo
			_, fuel := tx.FlightTime(f.origin.GetCoordinate(), f.destination, f.speed, f.ships)
			techs := tx.GetResearch()
			cargoCapacity := f.ships.Cargo(techs, f.b.GetServer().Settings.EspionageProbeRaids == 1)
			planetResources, _ := tx.GetResources(f.origin.GetID())
			planetResources.Deuterium -= fuel + 10
			payload.Deuterium = int(math.Min(float64(cargoCapacity), float64(planetResources.Deuterium)))
			cargoCapacity -= payload.Deuterium
			payload.Crystal = int(math.Min(float64(cargoCapacity), float64(planetResources.Crystal)))
			cargoCapacity -= payload.Crystal
			payload.Metal = int(math.Min(float64(cargoCapacity), float64(planetResources.Metal)))
		}

		f.fleet, f.err = tx.EnsureFleet(f.origin.GetID(), f.ships.ToQuantifiables(), f.speed, f.destination, f.mission, payload, f.expeditiontime, f.unionID)
		return f.err
	})
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
