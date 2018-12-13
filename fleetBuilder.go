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
	f := new(FleetBuilderFactory)
	f.b = b
	return f
}

// NewFleet ...
func (f FleetBuilderFactory) NewFleet() *FleetBuilder {
	return NewFleetBuilder(f.b)
}

// FleetBuilder ...
type FleetBuilder struct {
	b                Wrapper
	origin           CelestialID
	destination      Coordinate
	speed            Speed
	mission          MissionID
	ships            []Quantifiable
	resources        Resources
	err              error
	fleet            Fleet
	expeditiontime   int
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
	c := f.b.GetCachedCelestial(v)
	if c != nil {
		f.origin = c.GetID()
	}
	return f
}

// SetDestination ...
func (f *FleetBuilder) SetDestination(v interface{}) *FleetBuilder {
	c := f.b.GetCachedCelestial(v)
	if c != nil {
		f.destination = c.GetCoordinate()
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

// AddShips ...
func (f *FleetBuilder) AddShips(id ID, nbr int) *FleetBuilder {
	if id.IsShip() && id != SolarSatelliteID && nbr > 0 {
		f.ships = append(f.ships, Quantifiable{id, nbr})
	}
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

// SendNow send the fleet with defined configurations
func (f *FleetBuilder) SendNow() (Fleet, error) {
	err := f.b.Tx(func(tx *Prioritize) error {

		// Set all ships
		if f.allShips {
			ships, _ := tx.GetShips(f.origin)
			f.ships = make([]Quantifiable, 0)
			for _, ship := range Ships {
				if ship.GetID() == SolarSatelliteID {
					continue
				}
				nbr := ships.ByID(ship.GetID())
				if nbr > 0 {
					f.ships = append(f.ships, Quantifiable{ship.GetID(), nbr})
				}
			}
		}

		// Calculate cargo
		cargoCapacity := 0
		for _, ship := range f.ships {
			cargoCapacity += Objs.ByID(ship.ID).(Ship).GetCargoCapacity() * ship.Nbr
		}

		payload := f.resources
		// Send all resources
		if f.resources.Metal == -1 && f.resources.Crystal == -1 && f.resources.Deuterium == -1 {
			planetResources, _ := tx.GetResources(f.origin)
			payload.Deuterium = int(math.Min(float64(cargoCapacity), float64(planetResources.Deuterium)))
			cargoCapacity -= payload.Deuterium
			payload.Crystal = int(math.Min(float64(cargoCapacity), float64(planetResources.Crystal)))
			cargoCapacity -= payload.Crystal
			payload.Metal = int(math.Min(float64(cargoCapacity), float64(planetResources.Metal)))
		}

		f.fleet, f.err = tx.SendFleet(f.origin, f.ships, f.speed, f.destination, f.mission, payload, f.expeditiontime)
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
				f.b.CancelFleet(f.fleet.ID)
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
