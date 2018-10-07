package ogame

type MoonID int

// Moon ogame moon object
type Moon struct {
	ID         MoonID
	Img        string
	Name       string
	Diameter   int
	Coordinate Coordinate
	Fields     Fields
}
