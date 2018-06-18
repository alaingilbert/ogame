package buildable

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings"
	"github.com/alaingilbert/ogame/defenses"
	"github.com/alaingilbert/ogame/ships"
	"github.com/alaingilbert/ogame/technologies"
)

// Buildable ...
type Buildable interface {
	GetOGameID() ogame.ID
	IsAvailable(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches, int) bool
	GetPrice(int) ogame.Resources
	GetRequirements() map[ogame.ID]int
}

// GetByID ...
func GetByID(id ogame.ID) Buildable {
	if id.IsBuilding() {
		return buildings.GetByID(id)
	}
	if id.IsTech() {
		return technologies.GetByID(id)
	}
	if id.IsDefense() {
		return defenses.GetByID(id)
	}
	if id.IsShip() {
		return ships.GetByID(id)
	}
	return nil
}
