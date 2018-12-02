package ogame

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
func (f *FleetBuilder) SetOrigin(id CelestialID) *FleetBuilder {
	f.origin = id
	return f
}

// SetDestination ...
func (f *FleetBuilder) SetDestination(destination Coordinate) *FleetBuilder {
	f.destination = destination
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

// SendNow send the fleet with defined configurations
func (f *FleetBuilder) SendNow() (Fleet, error) {
	err := f.b.Tx(func(tx *Prioritize) error {
		res := f.resources
		// Send all resources
		if f.resources.Metal == -1 && f.resources.Crystal == -1 && f.resources.Deuterium == -1 {
			res, _ = tx.GetResources(f.origin)
		}
		f.fleet, f.err = tx.SendFleet(f.origin, f.ships, f.speed, f.destination, f.mission, res, f.expeditiontime)
		return f.err
	})
	if err != nil {
		// On error, call error callbacks
		for _, clb := range f.errorCallbacks {
			clb(err)
		}
	} else {
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
