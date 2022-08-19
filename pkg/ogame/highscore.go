package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
	"strconv"
)

// Highscore ...
type Highscore struct {
	NbPage   int64
	CurrPage int64
	Category int64 // 1:Player, 2:Alliance
	Type     int64 // 0:Total, 1:Economy, 2:Research, 3:Military, 4:Military Built, 5:Military Destroyed, 6:Military Lost, 7:Honor
	Players  []HighscorePlayer
}

// String ...
func (h Highscore) String() string {
	return "" +
		"  NbPage: " + utils.FI64(h.NbPage) + "\n" +
		"CurrPage: " + utils.FI64(h.CurrPage) + "\n" +
		"Category: " + utils.FI64(h.Category) + "\n" +
		"    Type: " + utils.FI64(h.Type) + "\n" +
		" Players: " + strconv.Itoa(len(h.Players)) + "\n"
}

// HighscorePlayer ...
type HighscorePlayer struct {
	Position     int64
	ID           int64
	Name         string
	Score        int64
	AllianceID   int64
	HonourPoints int64
	Homeworld    Coordinate
	Ships        int64 // When getting military type
}

// String ...
func (h HighscorePlayer) String() string {
	return "" +
		"    Position: " + utils.FI64(h.Position) + "\n" +
		"          ID: " + utils.FI64(h.ID) + "\n" +
		"        Name: " + h.Name + "\n" +
		"       Score: " + utils.FI64(h.Score) + "\n" +
		"  AllianceID: " + utils.FI64(h.AllianceID) + "\n" +
		"HonourPoints: " + utils.FI64(h.HonourPoints) + "\n" +
		"   Homeworld: " + h.Homeworld.String() + "\n" +
		"       Ships: " + utils.FI64(h.Ships) + "\n"
}
