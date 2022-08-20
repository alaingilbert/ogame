package ogame

// Planet ogame planet object
type Planet struct {
	Img         string
	ID          PlanetID
	Name        string
	Diameter    int64
	Coordinate  Coordinate
	Fields      Fields
	Temperature Temperature
	Moon        *Moon
}

func (p Planet) GetID() CelestialID          { return p.ID.Celestial() }
func (p Planet) GetImg() string              { return p.Img }
func (p Planet) GetName() string             { return p.Name }
func (p Planet) GetDiameter() int64          { return p.Diameter }
func (p Planet) GetCoordinate() Coordinate   { return p.Coordinate }
func (p Planet) GetFields() Fields           { return p.Fields }
func (p Planet) GetTemperature() Temperature { return p.Temperature }
func (p Planet) GetMoon() *Moon              { return p.Moon }
func (p Planet) GetType() CelestialType      { return PlanetType }

type Moon struct {
	ID         MoonID
	Img        string
	Name       string
	Diameter   int64
	Coordinate Coordinate
	Fields     Fields
}

func (m Moon) GetID() CelestialID        { return m.ID.Celestial() }
func (m Moon) GetImg() string            { return m.Img }
func (m Moon) GetName() string           { return m.Name }
func (m Moon) GetDiameter() int64        { return m.Diameter }
func (m Moon) GetCoordinate() Coordinate { return m.Coordinate }
func (m Moon) GetFields() Fields         { return m.Fields }
func (p Moon) GetType() CelestialType    { return MoonType }
