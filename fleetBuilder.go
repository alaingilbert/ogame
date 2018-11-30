package ogame

// FleetBuilder ...
type FleetBuilder struct {
	b           Wrapper
	origin      CelestialID
	destination Coordinate
	speed       Speed
	mission     MissionID
	ships       []Quantifiable
	resources   Resources
	err         error
	fleet       Fleet
}

// NewFleetBuilder ...
func NewFleetBuilder() *FleetBuilder {
	b := new(FleetBuilder)
	b.mission = Transport
	return b
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

// SetAllResources ...
func (f *FleetBuilder) SetAllResources() *FleetBuilder {
	f.resources = Resources{Metal: -1, Crystal: -1, Deuterium: -1}
	return f
}

// SetMission ...
func (f *FleetBuilder) SetMission(mission MissionID) *FleetBuilder {
	f.mission = mission
	return f
}

// AddShips ...
func (f *FleetBuilder) AddShips(id ID, nbr int) *FleetBuilder {
	if id.IsShip() && id != SolarSatelliteID && nbr > 0 {
		f.ships = append(f.ships, Quantifiable{id, nbr})
	}
	return f
}

// SendNow ...
func (f *FleetBuilder) SendNow() *FleetBuilder {
	res := f.resources
	// Send all resources
	if f.resources.Metal == -1 && f.resources.Crystal == -1 && f.resources.Deuterium == -1 {
		res, _ = f.b.GetResources(f.origin)
	}
	f.fleet, f.err = f.b.SendFleet(f.origin, f.ships, f.speed, f.destination, f.mission, res)
	return f
}

// OnError ...
func (f *FleetBuilder) OnError(clb func(error)) *FleetBuilder {
	clb(f.err)
	return f
}

// OnSuccess ...
func (f *FleetBuilder) OnSuccess(clb func(Fleet)) *FleetBuilder {
	clb(f.fleet)
	return f
}
