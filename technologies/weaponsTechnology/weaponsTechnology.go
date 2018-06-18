package weaponsTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// WeaponsTechnology ...
type WeaponsTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *WeaponsTechnology {
	b := new(WeaponsTechnology)
	b.OGameID = 109
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 800, Crystal: 200}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 4}
	return b
}
