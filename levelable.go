package ogame

//import (
//	"github.com/alaingilbert/ogame"
//	"github.com/alaingilbert/ogame/buildings"
//	"github.com/alaingilbert/ogame/technologies"
//)
//
//// Levelable ...
//type Levelable interface {
//	Buildable
//	GetLevel(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches) int
//	ConstructionTime(level, universeSpeed int, facilities ogame.Facilities) int
//}
//
//// GetByID ...
//func GetByID(id ogame.ID) Levelable {
//	if id.IsBuilding() {
//		return buildings.GetByID(id)
//	} else if id.IsTech() {
//		return technologies.GetByID(id)
//	}
//	return nil
//}
