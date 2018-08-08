package ogame

//// Buildable ...
//type Buildable interface {
//	GetID() ID
//	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
//	GetPrice(int) Resources
//	GetRequirements() map[ID]int
//}
//
//// GetByID ...
//func GetByID(id ID) Buildable {
//	if id.IsBuilding() {
//		return buildings.GetByID(id)
//	}
//	if id.IsTech() {
//		return technologies.GetByID(id)
//	}
//	if id.IsDefense() {
//		return defenses.GetByID(id)
//	}
//	if id.IsShip() {
//		return ships.GetByID(id)
//	}
//	return nil
//}
