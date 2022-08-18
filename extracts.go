package ogame

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = ExtractorV6{}
var _ Extractor = (*ExtractorV6)(nil)
var _ Extractor = ExtractorV7{}
var _ Extractor = (*ExtractorV7)(nil)

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int64 {
	return extractUniverseSpeedV6(pageHTML)
}

// ExtractorPlanet ogame planet object
type ExtractorPlanet struct {
	img         string
	id          PlanetID
	name        string
	diameter    int64
	coordinate  Coordinate
	fields      Fields
	temperature Temperature
	moon        *ExtractorMoon
}

func (p ExtractorPlanet) CelestialID() CelestialID { return p.id.Celestial() }
func (p ExtractorPlanet) Img() string              { return p.img }
func (p ExtractorPlanet) ID() PlanetID             { return p.id }
func (p ExtractorPlanet) Name() string             { return p.name }
func (p ExtractorPlanet) Diameter() int64          { return p.diameter }
func (p ExtractorPlanet) Coordinate() Coordinate   { return p.coordinate }
func (p ExtractorPlanet) Fields() Fields           { return p.fields }
func (p ExtractorPlanet) Temperature() Temperature { return p.temperature }
func (p ExtractorPlanet) Moon() *ExtractorMoon     { return p.moon }

type ExtractorMoon struct {
	id         MoonID
	img        string
	name       string
	diameter   int64
	coordinate Coordinate
	fields     Fields
}

func (m ExtractorMoon) CelestialID() CelestialID { return m.id.Celestial() }
func (m ExtractorMoon) ID() MoonID               { return m.id }
func (m ExtractorMoon) Img() string              { return m.img }
func (m ExtractorMoon) Name() string             { return m.name }
func (m ExtractorMoon) Diameter() int64          { return m.diameter }
func (m ExtractorMoon) Coordinate() Coordinate   { return m.coordinate }
func (m ExtractorMoon) Fields() Fields           { return m.fields }
